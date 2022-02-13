export class Attachment {

	// parent id to specify if no parent should be assign to this attachment
	public static readonly PARENT_DEFAULT = 'GLOBAL';
	public static readonly TYPE_DEFAULT = 'gallery';
	public static readonly GROUP_DEFAULT = 'DEFAULT';

	private static EXTENSION_ICON: Array<{ value: string, icon: string }> = [
		{value: 'pdf', icon: 'picture_as_pdf'},
		{value: 'docx', icon: 'text_format'},
		{value: 'zip', icon: 'file_copy'},
		{value: 'csv', icon: 'alternate_email'},
		{value: 'jpg', icon: 'image'},
		{value: 'jpeg', icon: 'image'},
		{value: 'svg', icon: 'image'},
		{value: 'png', icon: 'image'},
		{value: 'mp4', icon: 'videocam'},
		{value: 'webm', icon: 'videocam'},
	];

	id: string;
	parentKey: string;
	parentType: string;
	name: string;
	group: string;
	type: string;
	description: string;
	resourceUrl: string;
	resourceThumbUrl: string;
	isMedia: boolean;
	altText: string;
	videoDurationSec: number;
	fileType: string;
	displayOrder: number;

	public static getIconForAttachmentName(name: string) {
		const iconDefault = 'file_copy';
		if (!name) {
			return iconDefault;
		}
		const extensions = name.split('.');
		if (extensions.length > 0) {
			const extension = extensions[extensions.length - 1];
			const support = this.EXTENSION_ICON.find((e) => {
				return e.value === extension;
			});
			if (support) {
				return support.icon;
			}
		}
		return iconDefault;
	}

	public isNew() {
		return !this.id;
	}

	constructor(json?: any) {
		if (json) {
			this.id = json.id;
			this.parentType = json.parentType;
			this.parentKey = json.parentKey;
			this.name = json.name;
			this.group = json.group;
			this.type = json.type;
			this.description = json.description;
			this.resourceUrl = json.resourceUrl;
			this.resourceThumbUrl = json.resourceThumbUrl;
			this.isMedia = json.isMedia;
			this.altText = json.altText;
			this.videoDurationSec = json.videoDurationSec;
			this.fileType = json.fileType;
			this.displayOrder = json.displayOrder;
		}
	}
}
