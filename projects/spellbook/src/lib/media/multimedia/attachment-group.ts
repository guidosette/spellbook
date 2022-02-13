/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-02-28.
 */
import {Attachment} from './attachment';

export class AttachmentGroup {

	attachments: Array<Attachment>;
	public maxItems: number; // if not set = 0 = infinite
	public description: string;

	constructor(public readonly name: string, public readonly type: string, public parentKey: string | number, maxItems?: number, description?: string) {
		this.maxItems = maxItems;
		this.description = description;
		this.attachments = new Array<Attachment>();
	}

	public addAttachment(att: Attachment): void {
		this.attachments.push(att);
	}
}
