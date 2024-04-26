package plugins

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"log"
	"net/url"
	"reflect"
	"testing"
)

func TestConvertToStorageDisk(t *testing.T) {
	type args struct {
		disk *Disk
	}
	sataDisk := &Disk{
		Name:  "",
		Bus:   "sata",
		File:  "sata.file",
		Index: 0,
	}

	sataStorage := &proxmox.QemuStorages{
		Sata: &proxmox.QemuSataDisks{
			Disk_0: &proxmox.QemuSataStorage{
				Passthrough: &proxmox.QemuSataPassthrough{File: sataDisk.File},
			},
		},
	}

	ideDisk := &Disk{
		Name:  "",
		Bus:   "ide",
		File:  "ide.fil",
		Index: 1,
	}

	ideStorage := &proxmox.QemuStorages{
		Ide: &proxmox.QemuIdeDisks{
			Disk_1: &proxmox.QemuIdeStorage{
				Passthrough: &proxmox.QemuIdePassthrough{File: ideDisk.File},
			},
		},
	}

	scsiDisk := &Disk{
		Name:  "",
		Bus:   "scsi",
		File:  "scsi.file",
		Index: 2,
	}

	scsiStorage := &proxmox.QemuStorages{
		Scsi: &proxmox.QemuScsiDisks{
			Disk_2: &proxmox.QemuScsiStorage{
				Passthrough: &proxmox.QemuScsiPassthrough{File: scsiDisk.File},
			},
		},
	}

	virtioDisk := &Disk{
		Name:  "",
		Bus:   "virtio",
		File:  "virtio.file",
		Index: 3,
	}

	virtioStorage := &proxmox.QemuStorages{
		VirtIO: &proxmox.QemuVirtIODisks{
			Disk_3: &proxmox.QemuVirtIOStorage{
				Passthrough: &proxmox.QemuVirtIOPassthrough{File: virtioDisk.File},
			},
		},
	}
	tests := []struct {
		name string
		args args
		want *proxmox.QemuStorages
	}{
		// TODO: Add test cases.
		{"sata-disk", args{disk: sataDisk}, sataStorage},
		{"ide-disk", args{disk: ideDisk}, ideStorage},
		{"scsi-disk", args{disk: scsiDisk}, scsiStorage},
		{"virtio-disk", args{disk: virtioDisk}, virtioStorage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToStorageDisk(tt.args.disk); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToStorageDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diskApply(t *testing.T) {
	type args struct {
		fv     interface{}
		target interface{}
		index  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := diskApply(tt.args.fv, tt.args.target, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("diskApply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateDomain(t *testing.T) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	api, err := url.JoinPath("https://", fmt.Sprintf("%s:%d", "192.168.8.50", 8006), "/api2/json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := proxmox.NewClient(api, nil, "", tlsConfig, "", 60)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Login("root@pam", "Howlink@1401", "")
	if err != nil {
		log.Fatal(err)
	}

	var qemuConfig proxmox.ConfigQemu
	qemuConfig.VmID = 1000
	qemuConfig.Memory = 2048
	qemuConfig.QemuCpu = "host"
	qemuConfig.QemuOs = "l26"
	qemuConfig.QemuVcpus = 2
	qemuConfig.QemuSockets = 2
	qemuConfig.QemuCores = 1
	qemuConfig.Name = "test"

	disk := &proxmox.QemuStorages{
		Sata: &proxmox.QemuSataDisks{
			Disk_0: &proxmox.QemuSataStorage{
				Passthrough: &proxmox.QemuSataPassthrough{
					File: "/mnt/pve/jvdqmi/images/000/6000C294-e027-9d31-7faa-e7202ab34af3-flat.vmdk",
				},
			},
		},
	}

	qemuConfig.Disks = disk

	net := make(map[string]interface{})
	net["bridge"] = "vmbr0"
	net["firewall"] = true
	net["model"] = "virtio"
	qemuConfig.QemuNetworks = make(map[int]map[string]interface{})
	qemuConfig.QemuNetworks[0] = net

	vm := proxmox.NewVmRef(1000)
	vm.SetNode("hl")
	err = qemuConfig.Create(vm, client)
	if err != nil {
		log.Fatal(err)
	}
}
