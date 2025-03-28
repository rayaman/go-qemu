package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/melbahja/goph"
	"github.com/rayaman/go-qemu/pkg/image"
	"github.com/rayaman/go-qemu/pkg/image/formats"
	"github.com/rayaman/go-qemu/pkg/types"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/crypto/ssh"
)

type Credentials struct {
	User           string `json:"user"`
	Pass           string `json:"pass"`
	MachineID      string `json:"machine_id"`
	PublicKeyPath  string `json:"public_key_path"`
	PrivateKeyPath string `json:"private_key_path"`
}

type Controller struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func (c *Controller) Stop() {
	c.cancel()
}

func KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}
		if n == name {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found")
}

func StartMachine(machine types.Machine) *Controller {
	ctrl := &Controller{}

	go func() {
		ctx, cfunc := context.WithCancel(context.TODO())
		ctrl.cancel = cfunc
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command(fmt.Sprintf("qemu-system-%v", machine.Arch), machine.Expand()...)
		fmt.Println(cmd.String())
		cmd.Stdout = &stderr
		cmd.Stderr = &stdout
		ctrl.cmd = cmd
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error:", err)
		}
		go func() {
			err := cmd.Wait()
			if err != nil {
				fmt.Println("Error:", err)
			}
			cfunc()
		}()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Killing Process")
				err := cmd.Process.Kill()
				if err != nil {
					fmt.Println("Error:", err)
				}
				return
			}
		}
	}()

	return ctrl
}

func main() {
	// var img image.Image = &formats.QCOW2{
	// 	ImageName:   "imgs/UbuntuTest.img",
	// 	BackingFile: "base/UbuntuBase.img",
	// 	BackingFmt:  "qcow2",
	// }

	// var img image.Image = &formats.QCOW2{
	// 	ImageName: "imgs/NetworkTest.img",
	// }

	// err := image.Create(img, types.GetSize(types.GB, 8), image.Options{IsBaseImage: true})
	// if err != nil {
	// 	panic(err)
	// }
	machine := types.Machine{
		Arch: types.X86_64,
		Cores: types.SMP{
			Cpus: 4,
		},
		Boot: types.Boot{
			Order: []types.Drives{types.HARDDISK},
			Menu:  types.Off,
		},
		Memory: types.Memory{
			Size: 4096,
		},
		Accel: types.Accel{
			Accelerator: types.WHPX,
		},
		Nic: types.NIC{
			Type: types.TAP,
			Options: &types.TapOptions{
				IFName: "qemu-tap",
			},
		},
		NoGraphic: types.Set,
		HardDiskA: `./imgs/UbuntuTest.img`,
	}
	ctrl := StartMachine(machine)

	fmt.Println("Press enter to stop machine")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	ctrl.Stop()
	fmt.Println("Press enter to exit")
	_, _ = reader.ReadString('\n')

	// Start new ssh connection with private key.

	//client, err := goph.New("baseimguser", "127.0.0.1", goph.Password("food99"))
	// _, err := SetUpMachine()
	// if err != nil {
	// 	panic(err)
	// }

	// auth, err := goph.Key("keys/3e6e04c3-1af4-41e7-8984-c60012034075", "")
	// if err != nil {
	// 	panic(err)
	// }
	// client, err := goph.NewConn(&goph.Config{
	// 	User:     "baseimguser",
	// 	Addr:     "127.0.0.1",
	// 	Port:     8888,
	// 	Auth:     auth,
	// 	Timeout:  goph.DefaultTimeout,
	// 	Callback: ssh.InsecureIgnoreHostKey(),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// defer client.Close()

	// out, err := client.Run(`echo 'baseimgpass' | sudo -S ls /tmp`)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(out))
	// SetUpMachine(types.GetSize(types.GB, 64))
}

func SetUpMachine(size types.Size) (*Credentials, error) {
	machine_id := uuid.New().String()
	private_path := filepath.Join("keys", machine_id)
	public_path := filepath.Join("keys", machine_id+".pub")
	err := MakeSSHKeyPair(public_path, private_path)

	if err != nil {
		return nil, err
	}

	img := &formats.QCOW2{
		ImageName:   fmt.Sprintf("imgs/%v.img", machine_id),
		BackingFile: "base/UbuntuBase.img",
		BackingFmt:  "qcow2",
	}

	err = image.Create(img, size)

	if err != nil {
		return nil, err
	}

	// ToDo start Machine

	Config := &Credentials{
		MachineID:      machine_id,
		PublicKeyPath:  public_path,
		PrivateKeyPath: private_path,
	}

	client, err := goph.NewConn(&goph.Config{
		User:     "baseimguser",
		Addr:     "127.0.0.1",
		Port:     8888,
		Auth:     goph.Password("baseimgpass"),
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return nil, err
	}

	defer client.Close()

	err = client.Upload(public_path, "/home/baseimguser/.ssh/authorized_keys")
	if err != nil {
		return nil, err
	}

	cmds := []string{
		// Allow us to use our private key to connect directly into root user
		`cp /home/baseimguser/.ssh/authorized_keys /root/.ssh/authorized_keys`,
	}

	for _, cmd := range cmds {
		out, err := client.Run(`echo 'baseimgpass' | sudo -S ` + cmd)

		if err != nil {
			return nil, fmt.Errorf("Got Error: %v --- %v", string(out), err)
		}
	}

	return Config, nil
}

func MakeSSHKeyPair(pubKeyPath, privateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	return os.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}
