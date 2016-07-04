package virt

import (
	"fmt"

	"github.com/alexzorin/libvirt-go"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
)

var _ env.Virt = (*Conn)(nil)

type Conn struct {
	libvirt.VirConnection
}

func NewConn() (*Conn, error) {
	c, err := libvirt.NewVirConnection("qemu:///system")
	if err != nil {
		return nil, err
	}
	return &Conn{VirConnection: c}, nil
}

func (c *Conn) StartVM(vm *model.VM) error {
	xml := fmt.Sprintf(`
		<domain type='qemu'>
			<name>%s</name>
			<memory unit='MiB'>%d</memory>
			<os>
				<type>hvm</type>
				<boot dev='network'/>
			</os>
			<devices>
				<graphics type='spice' autoport='yes'/>
				<interface type='direct'>
					<source dev='em1' mode='bridge'/>
				</interface>
			</devices>
		</domain>
		`, vm.Name, vm.MemSizeMB)
	libvirtVm, err := c.DomainCreateXML(xml, 0)

	if err != nil {
		return err
	}

	fmt.Println(libvirtVm.GetXMLDesc(0))

	return nil
}

func (c *Conn) StopVM(vm *model.VM) error {
	libvirtVm, err := c.LookupDomainByName(vm.Name)
	if err != nil {
		return err
	}
	return libvirtVm.Destroy()
}

func (c *Conn) ListAll() error {
	vms, err := c.ListAllDomains(0)
	if err != nil {
		return err
	}

	for _, vm := range vms {
		name, err := vm.GetName()
		if err != nil {
			return err
		}
		fmt.Println(name)
	}

	return nil
}
