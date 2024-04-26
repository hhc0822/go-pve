package plugins

import (
	"errors"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"reflect"
	"strconv"
	"strings"
)

type Disk struct {
	Name  string `json:"name"`
	Bus   string `json:"bus"`
	File  string `json:"file"`
	Index int32  `json:"index"`
}

func ConvertToStorageDisk(disk *Disk) *proxmox.QemuStorages {
	qs := &proxmox.QemuStorages{}
	index := int(disk.Index)
	switch disk.Bus {
	case "sata":
		sataDisk := &proxmox.QemuSataStorage{
			Passthrough: &proxmox.QemuSataPassthrough{File: disk.File},
		}
		targetDisk := &proxmox.QemuSataDisks{}
		diskApply(sataDisk, targetDisk, index)
		qs.Sata = targetDisk
	case "ide":
		ideDisk := &proxmox.QemuIdeStorage{
			Passthrough: &proxmox.QemuIdePassthrough{File: disk.File},
		}
		targetDisk := &proxmox.QemuIdeDisks{}
		diskApply(ideDisk, targetDisk, index)
		qs.Ide = targetDisk
	case "scsi":
		scsiDisk := &proxmox.QemuScsiStorage{
			Passthrough: &proxmox.QemuScsiPassthrough{File: disk.File},
		}
		targetDisk := &proxmox.QemuScsiDisks{}
		diskApply(scsiDisk, targetDisk, index)
		qs.Scsi = targetDisk
	case "virtio":
		virtioDisk := &proxmox.QemuVirtIOStorage{
			Passthrough: &proxmox.QemuVirtIOPassthrough{File: disk.File},
		}
		targetDisk := &proxmox.QemuVirtIODisks{}
		diskApply(virtioDisk, targetDisk, index)
		qs.VirtIO = targetDisk
	}

	return qs
}

func diskApply(fv, target interface{}, index int) error {

	fvv := reflect.ValueOf(fv)
	tv := reflect.ValueOf(target)
	tt := reflect.TypeOf(target)

	typ := tv.Type()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		tv = tv.Elem()
	}
	if tt.Kind() == reflect.Pointer {
		tt = tt.Elem()
	}
	if index > typ.NumField()-1 {
		return errors.New("index out of range")
	}
	field := tv.Field(index)
	if fvv.Type() == field.Type() {
		js := tt.Field(index).Tag.Get("json")
		tag := strings.TrimSuffix(js, ",omitempty")
		i, err := strconv.Atoi(tag)
		if err != nil {
			return err
		}
		if index == i {
			field.Set(fvv)
		}
	}
	return nil
}
