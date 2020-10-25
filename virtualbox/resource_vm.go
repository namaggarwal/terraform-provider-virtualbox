package virtualbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gvb "github.com/namaggarwal/go-virtualbox"
)

func resourceVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the resource",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource",
				ForceNew:    true,
			},
		},
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Delete: resourceVMDelete,
	}
}

func resourceVMCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(gvb.IVBoxManage)
	name := d.Get("name").(string)
	vm := gvb.VirtualMachine{
		Name:       name,
		BaseFolder: "/Users/naman/VirtualBox VMs",
	}
	uuid, err := client.CreateVM(vm, true)
	if err != nil {
		return err
	}
	d.SetId(uuid)
	d.Set("uuid", uuid)
	return nil
}

func resourceVMRead(d *schema.ResourceData, m interface{}) error {
	uuid := d.Id()
	client := m.(gvb.IVBoxManage)
	vm, err := client.VMInfo(uuid)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(vm.UUID)
	d.Set("uuid", vm.UUID)
	d.Set("name", vm.Name)
	return nil
}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
	uuid := d.Id()
	client := m.(gvb.IVBoxManage)
	err := client.UnRegisterVM(uuid, true)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
