package winrmcp

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/masterzen/winrm/winrm"
)

// Communicator represents the Upload work
type Upload struct {
	winrmcp  *Winrmcp
	filePath string
}

func (u *Upload) Start() error {
	// Do stuff
	// tempFile := fmt.Sprintf("winrmcp-%s.tmp", uuid.TimeOrderedUUID())
	// tempPath := "$env:TEMP\\" + tempFile
	//
	// if os.Getenv("WINRMCP_DEBUG") != "" {
	// 	log.Printf("Copying file to %s\n", tempPath)
	// }
	//
	// err := uploadContent(client, config.MaxOperationsPerShell, "%TEMP%\\"+tempFile, in)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Error uploading file to %s: %v", tempPath, err))
	// }
	//
	// if os.Getenv("WINRMCP_DEBUG") != "" {
	// 	log.Printf("Moving file from %s to %s", tempPath, toPath)
	// }
	//
	// err = restoreContent(client, tempPath, toPath)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Error restoring file from %s to %s: %v", tempPath, toPath, err))
	// }
	//
	// if os.Getenv("WINRMCP_DEBUG") != "" {
	// 	log.Printf("Removing temporary file %s", tempPath)
	// }
	//
	// err = cleanupContent(client, tempPath)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("Error removing temporary file %s: %v", tempPath, err))
	// }
	//
	// return nil
	var content string
	err := uploadContent(u.winrmcp.client, u.winrmcp.config.MaxOperationsPerShell, u.winrmcp.config.MaxShell, u.filePath, content)
	return err
}

func (u *Upload) Stop() error {
	// Do stuff
	var err error
	return err
}

func uploadContent(client *winrm.Client, maxChunks int, maxParallel int, filePath string, content string) error {
	var err error
	parallel := 4
	var wg sync.WaitGroup

	// if maxChunks == 0 {
	// 	maxChunks = 1
	// }
	//
	// // Create 4 Parallel workers
	// for p := 0; p < parallel; p++ {
	// 	done := make(chan bool, 1)
	// 	// Add worker to the WaitGroup
	// 	wg.Add(1)
	// 	var thread = p
	// 	go func() {
	// 		defer wg.Done()
	// 	Loop:
	// 		for {
	// 			select {
	// 			case <-done:
	// 				break Loop
	// 			default:
	// 				finished, err := uploadChunks(client, fmt.Sprintf("%v.%v", filePath, thread), maxChunks, reader, thread)
	// 				if err != nil {
	// 					break
	// 				}
	// 				if finished {
	// 					done <- true
	// 				}
	// 			}
	// 		}
	// 	}()
	// }
	// fmt.Println("Waiting for Threads")
	// wg.Wait()
	// fmt.Println("Done waiting for Threads")
	return err
}

func uploadChunks(client *winrm.Client, filePath string, maxChunks int, reader io.Reader, thread int) (bool, error) {
	var done bool

	shell, err := client.CreateShell()
	if err != nil {
		return false, fmt.Errorf("Couldn't create shell: %v", err)
	}
	defer shell.Close()

	// Each shell can do X amount of chunks per session
	for c := 0; c < maxChunks; c++ {
		// Read a chunk
		content, finished, err := getChunk(reader, filePath)
		if err != nil {
			return false, err
		}
		if finished {
			done = true
		} else {
			// Upload chunk
			err = appendContent(shell, filePath, content)
			if err != nil {
				return false, err
			}
		}
	}

	return done, err
}

func appendContent(shell *winrm.Shell, filePath, content string) error {
	cmd, err := shell.Execute(fmt.Sprintf("echo %s >> \"%s\"", content, filePath))

	if err != nil {
		return err
	}

	defer cmd.Close()
	go io.Copy(os.Stdout, cmd.Stdout)
	go io.Copy(os.Stderr, cmd.Stderr)
	cmd.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("upload operation returned code=%d", cmd.ExitCode())
	}

	return nil
}

func restoreContent(client *winrm.Client, fromPath, toPath string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}

	defer shell.Close()
	script := fmt.Sprintf(`
		$tmp_file_path = [System.IO.Path]::GetFullPath("%s")
		$dest_file_path = [System.IO.Path]::GetFullPath("%s")
		if (Test-Path $dest_file_path) {
			rm $dest_file_path
		}
		else {
			$dest_dir = ([System.IO.Path]::GetDirectoryName($dest_file_path))
			New-Item -ItemType directory -Force -ErrorAction SilentlyContinue -Path $dest_dir | Out-Null
		}

		if (Test-Path $tmp_file_path) {
			$base64_lines = Get-Content $tmp_file_path
			$base64_string = [string]::join("",$base64_lines)
			$bytes = [System.Convert]::FromBase64String($base64_string)
			[System.IO.File]::WriteAllBytes($dest_file_path, $bytes)
		} else {
			echo $null > $dest_file_path
		}
	`, fromPath, toPath)

	cmd, err := shell.Execute(winrm.Powershell(script))
	if err != nil {
		return err
	}

	go io.Copy(os.Stdout, cmd.Stdout)
	go io.Copy(os.Stderr, cmd.Stderr)

	cmd.Wait()
	cmd.Close()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("restore operation returned code=%d", cmd.ExitCode())
	}
	return nil
}

func cleanupContent(client *winrm.Client, filePath string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}

	defer shell.Close()
	cmd, _ := shell.Execute("powershell", "Remove-Item", filePath, "-ErrorAction SilentlyContinue")

	cmd.Wait()
	cmd.Close()
	return nil
}
