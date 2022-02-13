/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-08-24.
 */
import {Definition, Field} from '../core/typedef.service';
import {Content} from './content';
import {FormGroup} from '@angular/forms';
import {Client} from '../core/client';

export class ContentDefinition implements Definition<Content> {

	protected mfields: Array<string>; // list mandatory
	protected columnfields: Array<string>; // list column
	protected columnFieldOrder: string;
	protected columnOrder: string;
	protected intermediateSlugUrl: string;

	public readonly title: Field;
	public readonly isPublished: Field;

	public published: Field;
	public locale: Field;
	public slug: Field;
	public order: Field;
	public subtitle: Field;
	public topic: Field;
	public description: Field;
	public parent: Field;
	public body: Field;
	public code: Field;
	public editor: Field;
	public attachments: Field;
	public tags: Field;
	public cover: Field;

	public startDate: Field;
	public endDate: Field;

	public created: Field;
	public delete: Field;

	constructor(private readonly type: string, private readonly label: string, private readonly icon: string) {
		// build default fields
		this.title = new Field('title', 'Title');
		this.locale = new Field('locale', 'Language');
		this.published = new Field('published', 'Publication Date');
		this.isPublished = new Field('isPublished', 'Published');
		this.slug = new Field('slug', 'Link');
		this.order = new Field('order', 'Order');
		this.tags = new Field('tags', 'Tags');
		this.cover = new Field('cover', 'Cover image');
		this.subtitle = new Field('subtitle', 'Subtitle');
		this.topic = new Field('topic', 'Topic');
		this.description = new Field('description', 'Description');
		this.parent = new Field('parent', 'Parent');
		this.body = new Field('body', 'Body');
		this.code = new Field('code', 'Code');
		this.editor = new Field('editor', 'Editor');
		this.attachments = new Field('attachments', 'Attachments');
		this.created = new Field('created', 'Created');
		this.delete = new Field('delete', 'Delete');
		this.mfields = [
			this.title.id,
			this.locale.id,
			this.isPublished.id,
		];
		this.columnfields = [
			this.title.id,
			this.locale.id,
			this.code.id,
			this.published.id,
			this.order.id,
			this.isPublished.id,
			this.created.id,
			this.delete.id,
		];
		this.columnFieldOrder = 'Order';
		this.columnOrder = Client.orderAscKey;
		this.intermediateSlugUrl = 'news';
	}

	definitionFor(): string {
		return this.type;
	}

	field(id: string): Field {
		return this[id];
	}

	menuLabel(): string {
		return this.label;
	}

	menuIcon(): string {
		return this.icon;
	}

	public isMandatory(field: string) {
		return this.mandatoryFields().indexOf(field) > -1;
	}

	public mandatoryFields(): Array<string> {
		return this.mfields;
	}

	public columnFields(): Array<string> {
		return this.columnfields;
	}

	public isColumnField(field: string) {
		return this.columnFields().indexOf(field) > -1;
	}

	public listOrderField(): string {
		return this.columnFieldOrder;
	}

	public listOrder(): string {
		return this.columnOrder;
	}

	public getIntermediateSlugUrl(): string {
		return this.intermediateSlugUrl;
	}

	public beforeSend(content: Content) {}

	public addValidationRules(form: FormGroup) {}

}
