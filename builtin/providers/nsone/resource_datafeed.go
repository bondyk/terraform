package nsone

import (
	"github.com/hashicorp/terraform/helper/schema"

	nsone "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func dataFeedResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Create: DataFeedCreate,
		Read:   DataFeedRead,
		Update: DataFeedUpdate,
		Delete: DataFeedDelete,
	}
}

func dataFeedToResourceData(d *schema.ResourceData, df *data.Feed) {
	d.SetId(df.ID)
	d.Set("name", df.Name)
	d.Set("config", df.Config)
}

func resourceDataToDataFeed(d *schema.ResourceData) *data.Feed {
	config := make(data.Config)
	for k, v := range d.Get("config").(map[string]interface{}) {
		config[k] = v.(string)
	}
	return &data.Feed{
		Name:     d.Get("name").(string),
		Config:   config,
		SourceID: d.Get("source_id").(string),
	}
}

// DataFeedCreate creates an ns1 datafeed
func DataFeedCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.Client)
	df := resourceDataToDataFeed(d)
	if _, err := client.DataFeeds.Create(d.Get("source_id").(string), df); err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}

// DataFeedRead reads the datafeed for the given ID from ns1
func DataFeedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.Client)
	df, _, err := client.DataFeeds.Get(d.Get("source_id").(string), d.Id())
	if err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}

// DataFeedDelete delets the given datafeed from ns1
func DataFeedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.Client)
	_, err := client.DataFeeds.Delete(d.Get("source_id").(string), d.Id())
	d.SetId("")
	return err
}

// DataFeedUpdate updates the given datafeed with modified parameters
func DataFeedUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.Client)
	df := resourceDataToDataFeed(d)
	df.ID = d.Id()
	if _, err := client.DataFeeds.Update(d.Get("source_id").(string), df); err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}
