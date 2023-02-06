// resource_server.go
package main

import (
	/* "log"
	"net/http" */
	"bufio"
    "errors"
	"context"
	"os"
	"strconv"

	// "time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	filePath = "c:\\xdump\\resource.txt"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"uuid_content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uuid_content := d.Get("uuid_content").(string)
	// time := time.Now().String()
	var diags diag.Diagnostics

	file, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		return diag.FromErr(err)
		// return err
	}

	_, err1 := file.WriteString(uuid_content)


	if err1 != nil {
		return diag.FromErr(err)
		// return err
	}

    contentLen := len(uuid_content)
    // proposedLen := strconv.Itoa(d.Get("get_len").(int))
    currentID := d.Id()
    proposedLen,err := strconv.Atoi(currentID)

	if err != nil {
		return diag.FromErr(err)
		// return err
	}

    if contentLen != proposedLen {
        return diag.FromErr(errors.New("Len_Not_Matched " + strconv.Itoa(contentLen) + ":" + strconv.Itoa(proposedLen)))
    }

	d.SetId(strconv.Itoa(proposedLen))

	if err != nil {
		return diag.FromErr(err)
	}

	// return diags
	// return resourceServerRead(ctx, d, m)
    return diags
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// uuid_content := d.Get("uuid_content").(string)
	var diags diag.Diagnostics

	// id := d.Id()

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
    if err != nil {
       return diag.FromErr(err)
       // return err
    }

    fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)

    var content string

    for fileScanner.Scan(){
       content += fileScanner.Text()
    }

    contentLen := len(content)
    currentID := d.Id()
    proposedLen := len(currentID)


    if contentLen != proposedLen {
        return diag.FromErr(errors.New("Len_Not_Matched " + strconv.Itoa(contentLen) + ":" + strconv.Itoa(proposedLen)))
    }


    if err := d.Set("uuid_content", content); err != nil {
       // return err
       return diag.FromErr(err)
    }
	// return nil
	return diags
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	uuid_content := d.Get("uuid_content").(string)

	d.SetId(uuid_content)

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		return diag.FromErr(err)
		// return err
	}

	_, err1 := file.WriteString(uuid_content)
	/* time := time.Now().String()
	   _, err2 := file.WriteString("Appending uuid (" + uuid_count + ") >>> " + time +  "\n") */

	if err1 != nil {
		return diag.FromErr(err)
		// return err
	}

	return diags
	// return resourceServerRead(ctx, d, m)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")

	err := os.Remove(filePath)
	if err != nil {
		return diag.FromErr(err)
		// return err
	}

	return diags
	// return nil
}
