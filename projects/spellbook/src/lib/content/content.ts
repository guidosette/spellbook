import {Attachment} from '../media/multimedia/attachment';
import {AttachmentGroup} from '../media/multimedia/attachment-group';
import {Spellbook} from '../core/spellbook';

export class Content {
	id: string;
	idTranslate: string;
	parent: string;
	code: string;
	type: string;
	slug: string;
	title: string;
	subtitle: string;
	body: string;
	tags: Array<string>;
	category: string;
	topic: string;
	locale: string;
	revision: string;
	cover: string;
	order: number;
	description: string;
	editor: string;
	author: string;
	created: string;
	updated: string;
	published: string;
	isPublished: boolean;

	attachmentGroups: Array<AttachmentGroup>;

	// CONTENT_TYPE_EVENT
	startDate: string;
	endDate: string;
	hasStartDate: boolean;
	hasEndDate: boolean;

	getType(): string {
		return this.type;
	}

	public isNew() {
		return !this.id;
	}

	constructor(json?: any, spellbook?: Spellbook) {
		this.tags = new Array<string>();
		this.attachmentGroups = new Array<AttachmentGroup>();

		if (json) {
			this.id = json.id;
			this.idTranslate = json.idTranslate;
			this.parent = json.parent;
			this.type = json.type;
			this.slug = json.slug;
			this.code = json.code;
			this.title = json.title;
			this.subtitle = json.subtitle;
			this.description = json.description;
			this.body = json.body;
			if (json.tags) {
				for (const p of json.tags) {
					this.tags.push(p);
				}
			}

			this.cover = json.cover;
			this.category = json.category;
			this.topic = json.topic;
			this.locale = json.locale;
			this.revision = json.revision;
			this.author = json.author;
			this.editor = json.editor;
			this.created = json.created;
			this.updated = json.updated;
			this.published = json.published;
			this.isPublished = json.isPublished;
			this.order = json.order;

			// CONTENT_TYPE_EVENT
			this.startDate = json.startDate;
			this.endDate = json.endDate;
			this.hasStartDate = json.hasStartDate;
			this.hasEndDate = json.hasEndDate;

			// unpack attachments
			if (json.attachments) {
				for (const attachment of json.attachments) {
					const att = new Attachment(attachment);
					const group = this.getAttachmentGroup(att.group, att.type);
					group.addAttachment(att);
				}
			}

			// default attachments group
			if (spellbook) {
				const cat = spellbook.getSupportedCategories().find((s) => {
					return s.name === this.category;
				});
				if (cat !== undefined) {
					for (const defaultAttachmentGroup of cat.defaultAttachmentGroups) {
						this.checkAttachmentGroup(defaultAttachmentGroup.Name, defaultAttachmentGroup.Type, defaultAttachmentGroup.MaxItem, defaultAttachmentGroup.Description);
					}
				}
			}
		}
	}

	public checkAttachmentGroup(name: string, type: string, maxItems?: number, description?: string) {
		const index = this.attachmentGroups.findIndex((ag: AttachmentGroup) => {
			return ag.name === name;
		});
		if (index === -1) {
			const att: AttachmentGroup = new AttachmentGroup(name, type, this.id, maxItems, description);
			this.attachmentGroups.push(att);
		} else {
			this.attachmentGroups[index].maxItems = maxItems;
			this.attachmentGroups[index].description = description;
		}
	}

	public getAttachmentGroup(group: string, type: string) {
		if (!this.id) {
			return null;
		}

		let attachmentGroup = this.attachmentGroups.find(x => x.type === type && x.name === group);
		if (attachmentGroup) {
			return attachmentGroup;
		}
		attachmentGroup = new AttachmentGroup(group, type, this.id);
		this.attachmentGroups.push(attachmentGroup);
		return attachmentGroup;
	}

	public toJSON(): any {
		const attachments = new Array<Attachment>();
		for (const ag of this.attachmentGroups) {
			for (const att of ag.attachments) {
				attachments.push(att);
			}
		}

		return {
			id: this.id,
			idTranslate: this.idTranslate,
			parent: this.parent,
			type: this.type,
			slug: this.slug,
			title: this.title,
			subtitle: this.subtitle,
			description: this.description,
			body: this.body,
			tags: this.tags,
			category: this.category,
			topic: this.topic,
			locale: this.locale,
			revision: this.revision,
			cover: this.cover,
			author: this.author,
			code: this.code,
			editor: this.editor,
			created: this.created,
			updated: this.updated,
			published: this.published,
			isPublished: this.isPublished,
			order: this.order,
			attachments,
			// CONTENT_TYPE_EVENT
			startDate: this.startDate,
			endDate: this.endDate,
			hasStartDate: this.hasStartDate,
			hasEndDate: this.hasEndDate,
		};
	}
}

