export class Subscription {

	key: string;
	email: string;
	country: string;
	firstName: string;
	lastName: string;
	organization: string;
	created: string;
	updated: string;
	position: string;
	notes: string;

	public isNew() {
		return !(this.key && this.key !== undefined && this.key.length > 0);
	}

	constructor(json?: any) {

		if (json) {
			this.key = json.key;
			this.email = json.email;
			this.country = json.country;
			this.firstName = json.firstName;
			this.lastName = json.lastName;
			this.organization = json.organization;
			this.created = json.created;
			this.updated = json.updated;
			this.position = json.position;
			this.notes = json.notes;
		}
	}

	public toJSON(): any {
		return {
			key: this.key,
			email: this.email
		};
	}
}

