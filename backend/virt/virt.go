package virt

import (
	"fmt"

	"github.com/alexzorin/libvirt-go"
)

type Conn interface {
	StartVM() error
}

type conn struct {
	libvirt.VirConnection
}

func NewConn() (Conn, error) {
	c, err := libvirt.NewVirConnection("qemu:///system")
	if err != nil {
		return nil, err
	}
	return &conn{VirConnection: c}, nil
}

func (c *conn) StartVM() error {
	vm, err := c.DomainCreateXML(`
			<domain type='qemu'>
				<name>go-vm</name>
				<memory unit='MiB'>512</memory>
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
			`, 0)

	if err != nil {
		return err
	}

	fmt.Println(vm.GetXMLDesc(0))

	return nil
}

func (c *conn) ListAll() error {
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
