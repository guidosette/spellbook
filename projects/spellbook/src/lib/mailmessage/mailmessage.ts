export class MailMessage {

	id: number;
	recipient: string;
	sender: string;
	object: string;
	body: string;
	created: string;

	constructor(json?: any) {

		if (json) {
			this.id = json.id;
			this.recipient = json.recipient;
			this.sender = json.sender;
			this.object = json.object;
			this.body = json.body;
			this.created = json.created;
		}
	}

	public toJSON(): any {
		return {
			id: this.id,
			email: this.recipient
		};
	}
}

