package content

import (
	"context"
	"distudio.com/mage/model"
	"distudio.com/page/resource"
	"google.golang.org/appengine/log"
)

type Manager struct {}

func (manager Manager) NewResource(ctx context.Context) (resource.Resource, error) {
	return &Content{}, nil
}


func (manager Manager) FromId(ctx context.Context, id string) (resource.Resource, error) {
	cont := Content{}
	if err := model.FromStringID(ctx, &cont, id, nil); err != nil {
		log.Errorf(ctx, "could not retrieve content %s: %s", id, err.Error())
		return nil, err
	}

	q := model.NewQuery((*Attachment)(nil))
	q = q.WithField("Parent =", cont.Slug)
	if err := q.GetMulti(ctx, &cont.Attachments); err != nil {
		log.Errorf(ctx, "could not retrieve content %s attachments: %s", id, err.Error())
		return nil, err
	}

	return &cont, nil
}

func (manager Manager) ListOf(ctx context.Context, opts resource.ListOptions) ([]resource.Resource, error) {
	var conts []*Content
	q := model.NewQuery(&Content{})
	q = q.OffsetBy(opts.Page * opts.Size)

	if opts.Order != "" {
		dir := model.ASC
		if opts.Descending {
			dir = model.DESC
		}
		q = q.OrderBy(opts.Order, dir)
	}

	// get one more so we know if we are done
	q = q.Limit(opts.Size + 1)
	err := q.GetMulti(ctx, &conts)
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, len(conts))
	for i := range conts {
		resources[i] = resource.Resource(conts[i])
	}

	return resources, nil
}

func (manager Manager) Save(ctx context.Context, res resource.Resource) error {
	content := res.(*Content)
	// input is valid, create the resource
	opts := model.CreateOptions{}
	opts.WithStringId(content.Slug)

	// // WARNING: the volatile field Multimedia because Memcache (Gob)
	//	can't ignore field
	tmp := content.Attachments
	content.Attachments = nil

	err := model.CreateWithOptions(ctx, content, &opts)
	if err != nil {
		log.Errorf(ctx, "error creating post %s: %s", content.Slug, err)
		return err
	}

	// return the swapped multimedia value
	content.Attachments = tmp
	return nil
}

func (manager Manager) Delete(ctx context.Context, resource resource.Resource) error {
	content := resource.(*Content)
	err := model.Delete(ctx, content, nil)
	if err != nil {
		log.Errorf(ctx, "error deleting content %s: %s", content.Slug, err.Error())
		return err
	}

	// delete attachments with parent = slug
	attachments := make([]*Attachment, 0, 0)
	q := model.NewQuery(&Attachment{})
	q.WithField("Parent =", content.Slug)
	err = q.GetMulti(ctx, &attachments)
	if err != nil {
		log.Errorf(ctx, "error retrieving attachments: %s", err)
		return err
	}

	for _, attachment := range attachments {
		err = model.Delete(ctx, attachment, nil)
		if err != nil {
			log.Errorf(ctx,"error deleting attachment %+v: %s", attachment, err.Error())
			return err
		}
	}

	return nil
}

func (manager Manager) PropertyValues(ctx context.Context, properties []string) ([]string, error) {

}

