export class Category {

	name: string;
	label: string;
	type: string;
	defaultAttachmentGroups: Array<DefaultAttachmentGroup>;

	constructor(json?: any) {
		if (json) {
			this.name = json.name;
			this.label = json.label;
			this.type = json.type;

			this.defaultAttachmentGroups = new Array<DefaultAttachmentGroup>();
			if (json.defaultAttachmentGroups) {
				for (const jsonDefault of json.defaultAttachmentGroups) {
					const defaultAttachmentGroup = new DefaultAttachmentGroup(jsonDefault);
					this.defaultAttachmentGroups.push(defaultAttachmentGroup);
				}
			}
		}
	}

}

export class DefaultAttachmentGroup {

	Name: string;
	Type: string;
	MaxItem: number;
	Description: string;

	constructor(json?: any) {
		if (json) {
			this.Name = json.Name;
			this.Type = json.Type;
			this.MaxItem = json.MaxItem;
			this.Description = json.Description;
		}
	}

}
