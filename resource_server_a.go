// resource_server.go
package main

import (
	/* "log"
	"net/http" */
	// "bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"log"
	// "strconv"

	// "errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	// "strconv"

	// "time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	fileRoot = "c:\\xdump"
    logFile *os.File
    logger *log.Logger
    logFilePath = "c:\\xdump\\log.txt"
)


func lo(message string){
    logger.Output(2, message)
}

func resourceServerA() *schema.Resource {

    _ = os.Remove(logFilePath)
    var err error
    logFile, err = os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE, 0644)
    if err != nil {
        log.Fatal(err)
    }
    logger = log.New(logFile, "",0)
    // logger.Output(2, "begin - resourceserverA")
    lo("begin resourceSeverA - new")

	return &schema.Resource{
		CreateContext: resourceServerCreateA,
		ReadContext:   resourceServerReadA,
		UpdateContext: resourceServerUpdateA,
		DeleteContext: resourceServerDeleteA,

		Schema: map[string]*schema.Schema{
            "pod": &schema.Schema{
                Type: schema.TypeList,
                MaxItems: 1,
                Required: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": &schema.Schema{
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "build": &schema.Schema{
                            Type: schema.TypeList,
                            Required: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "ram": &schema.Schema{
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "disk": &schema.Schema{
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "size": &schema.Schema{
                                        Type:     schema.TypeInt,
                                        Computed: true,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
	}
}

type PodCon struct {
    Name string
    Build Build
}
type Build struct {
    Ram int
    Disk int
}

func getHash (filePath string) (string, error) {
	file, err := os.Open(filePath) // Open the file for reading
	if err != nil {
        return "", err
	}
	defer file.Close() // Be sure to close your file!


	hash := sha256.New() // Use the Hash in crypto/sha256

	if _, err := io.Copy(hash, file); err != nil {
        return "", err
	}

	sum := fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
    return sum, nil
}


func resourceServerCreateA(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lo("resourceCreate A")
    var diags diag.Diagnostics
    var podList []PodCon

    pods := d.Get("pod").([]interface{})
    lo("get props ")

    for _, pod := range pods {
        var podObj PodCon
        var buildObj Build
        p := pod.(map[string]interface{})
        bld := p["build"].([]interface{})[0]
        build := bld.(map[string]interface{})

        buildObj.Ram = build["ram"].(int)
        buildObj.Disk = build["disk"].(int)

        podObj.Name = p["name"].(string)
        podObj.Build = buildObj

        podList = append(podList, podObj)
    }

    lo("server config list collected from d interface")


    for _, podO := range podList {
        name := podO.Name
        filePath = fileRoot + "\\" + name + ".json"
        _ = os.Remove(filePath)
        file, err := os.OpenFile(filePath, os.O_CREATE , 0644)
        if err != nil {
            return diag.FromErr(err)
            // return err
        }
        defer file.Close()
        serverObject, _ := json.MarshalIndent(podO,"","  ")

        _, err1 := file.Write(serverObject)
        if err1 != nil {
            // return err
        }
    }

    hash, err := getHash(filePath)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash)
    resourceServerReadA(ctx, d, m)
    return diags
}

func resourceServerReadA(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    files, err := ioutil.ReadDir(fileRoot)
    if err != nil {
        log.Fatal(err)
    }

    lo("in read")
    hash := d.Id()
    lo("current d.id hash is " + hash)

    for _, file := range files {
        if !file.IsDir(){

            filePath := fileRoot + "\\" + file.Name()
            fHash, err := getHash(filePath)
            if err != nil {
                return diag.FromErr(err)
            }
            lo("file hash is " + fHash)
            if hash == fHash {
                var podC PodCon
                b, err := ioutil.ReadFile(filePath)
                if err != nil {
                    return diag.FromErr(err)
                }
                err1 := json.Unmarshal(b, &podC)
                if err1 != nil {
                    return diag.FromErr(err)
                }
                result := flattenPod (podC)
                if err := d.Set("pod", result); err != nil {
                    return diag.FromErr(err)
                }
                break
            } else {
                lo("Setting new hash")
                d.SetId("")
            }
        }
    }

    return diags

}

func resourceServerUpdateA (ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var podList []PodCon

    lo("in update")
    if d.HasChange("pod") {

        pods := d.Get("pod").([]interface{})
        for _, pod := range pods {
            var podObj PodCon
            var buildObj Build
            p := pod.(map[string]interface{})
            bld := p["build"].([]interface{})[0]
            build := bld.(map[string]interface{})

            buildObj.Ram = build["ram"].(int)
            buildObj.Disk = build["disk"].(int)

            podObj.Name = p["name"].(string)
            podObj.Build = buildObj

            podList = append(podList, podObj)
        }

        lo("server config list collected from d interface")


        for _, podO := range podList {
            name := podO.Name
            filePath = fileRoot + "\\" + name + ".json"
            file, err := os.OpenFile(filePath, os.O_CREATE , 0644)
            if err != nil {
                return diag.FromErr(err)
                // return err
            }
            defer file.Close()
            serverObject, _ := json.MarshalIndent(podO,"","  ")

            _, err1 := file.Write(serverObject)
            if err1 != nil {
                // return err
            }
            hash, err := getHash(filePath)

            if err != nil {
                return diag.FromErr(err)
            }

            d.SetId(hash)
        }

    }
    return resourceServerReadA(ctx, d, m)
}

func resourceServerDeleteA(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func flattenPod(pod PodCon) []interface{} {
		ois := make([]interface{}, 1,1)

        oi := make(map[string]interface{})

        oi["build"] = flattenBuild(pod.Build)
        oi["name"] = pod.Name
        ois[0] = oi

		return ois

	return make([]interface{}, 0)
}

func flattenBuild(build Build) []interface{} {
	c := make(map[string]interface{})
	c["ram"] = build.Ram
	c["disk"] = build.Disk
	return []interface{}{c}
}
