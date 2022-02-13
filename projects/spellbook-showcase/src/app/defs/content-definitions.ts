/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-08-25.
 */
import {ContentDefinition} from '../../../../spellbook/src/lib/content/content-definition';
import {Field} from '../../../../spellbook/src/lib/core/typedef.service';
import {Client} from '../../../../spellbook/src/lib/core/client';


export class SpecialDefinition extends ContentDefinition {

	constructor() {
		super('example_type', 'Example type', 'folder_special');

		this.tags = null;
		this.topic = null;
		this.description.label = 'Description';
		this.subtitle = null;
		this.code.label = 'Code';
		// this.attachments = null;
		this.parent = null;
		this.body.label = 'Body';
		this.editor = null;
		this.mfields = [
			this.title.id,
			this.description.id,
			this.isPublished.id,
			this.code.id,
		];
		this.columnfields = [
			this.title.id,
			this.locale.id,
			this.code.id,
			this.created.id,
			this.order.id,
			this.isPublished.id,
			this.delete.id,
		];
		this.columnFieldOrder = 'Title';
		this.columnOrder = Client.orderDescKey;
		this.intermediateSlugUrl = 'news3';
	}
}

export class EventDefinition extends ContentDefinition {

	constructor() {
		super('events_type', 'Events type', 'event');
		this.locale.hidden = true;
		this.subtitle = new Field('subtitle', 'Excerpt');
		this.endDate = new Field('endDate', 'Data fine');
		this.startDate = new Field('startDate', 'Data inizio');
		this.subtitle.label = 'Location';
		this.description.label = 'Excerpt';
		this.editor = null;
		this.parent = null;
		this.attachments = null;
		this.code = null;
		this.order = null;
		this.topic.label = 'type';
		this.mfields = [
			this.title.id,
			this.locale.id,
			this.isPublished.id,
			this.startDate.id,
		];
		this.columnfields = [
			this.slug.id,
			this.title.id,
			this.startDate.id,
			this.endDate.id,
			this.published.id,
			this.isPublished.id,
			this.created.id,
			this.delete.id,
		];
		this.columnFieldOrder = 'Title';
		this.columnOrder = Client.orderAscKey;
		this.intermediateSlugUrl = 'news2';
	}

}
