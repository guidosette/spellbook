export class Place {

	id: number;
	name: string;
	address: string;
	city: string;
	street: string;
	streetNumber: string;
	postalCode: string;
	area: string;
	country: string;
	description: string;
	phone: string;
	created: string;
	website: string;
	lat: number;
	lng: number;

	constructor(json?: any) {

		if (json) {
			this.id = json.id;
			this.name = json.name;
			this.address = json.address;
			this.street = json.street;
			this.area = json.area;
			this.streetNumber = json.streetNumber;
			this.city = json.city;
			this.postalCode = json.postalCode;
			this.country = json.country;
			this.description = json.description;
			this.phone = json.phone;
			this.created = json.created;
			this.lat = json.lat;
			this.lng = json.lng;
			this.website = json.website;
		}
	}

	public isNew() {
		return !(this.id && this.id !== 0);
	}

	public toJSON(): any {
		return {
			id: this.id,
			name: this.name,
			address: this.address,
			street: this.street,
			streetNumber: this.streetNumber,
			area: this.area,
			city: this.city,
			postalCode: this.postalCode,
			country: this.country,
			description: this.description,
			phone: this.phone,
			created: this.created,
			lat: this.lat,
			lng: this.lng,
			website: this.website,
		};
	}
}

