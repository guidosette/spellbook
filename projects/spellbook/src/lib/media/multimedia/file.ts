export class FileAttachment {

	// File
	name: string;
	resourceUrl: string;
	resourceThumbUrl: string;
	contentType: string; // ContentType is the MIME type of the object's content.

	attachmentType: string; // AttachmentType

	constructor(json?: any) {
		if (json) {
			this.name = json.name;
			this.resourceUrl = json.resourceUrl;
			this.resourceThumbUrl = json.resourceThumbUrl;
			this.contentType = json.contentType;
		}
	}
}
