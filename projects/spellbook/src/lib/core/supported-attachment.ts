
export class AttachmentType {
	public static readonly GALLERY = 'gallery';
	public static readonly ATTACHMENT = 'attachments';
	public static readonly VIDEO = 'video';
}

export class SupportedAttachment {

	name: string;
	value: string; // AttachmentType
	image: string;
	accept: string;
	mime: string;
}
