package main

import (
	"bufio"
	"context"
	"os"
	"strconv"
	// "time"
    "errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCoffees() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"get_len": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"uuid_content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	filePath := "c:\\xdump\\resource.txt"
	var diags diag.Diagnostics
    // var getLen string
    getLen := strconv.Itoa(d.Get("get_len").(int))
    // getLen := strconv.Itoa(d.Get("get_len").(int))
    // getLen := d.Get("get_len")
    proposedLen := len(getLen)

	file, err := os.OpenFile(filePath, os.O_APPEND, 0644)
	if err != nil {
		return diag.FromErr(err)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var content string

	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
    contentLen := len(content)
    if proposedLen != contentLen {
        return diag.FromErr(errors.New("Len_Not_Matched " + strconv.Itoa(contentLen) + ":" + strconv.Itoa(proposedLen)))
        // return diag.FromErr(errors.New("Len not match"))
        // return diag.FromErr(errors.New("Len_Not_Matched"))
    }

	if err := d.Set("uuid_content", content); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(getLen)
	// _, err1 := file.
	/* time := time.Now().String()
	   _, err2 := file.WriteString("Appending uuid (" + uuid_count + ") >>> " + time +  "\n") */

	/* if err1 != nil {
	    return diag.FromErr(err)
	} */

	return diags
}
