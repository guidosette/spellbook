export class Page {

	id: number;
	label: string;
	title: string;
	metadesc: string;
	url: string;
	order: number;
	locale: string;
	code: string;

	constructor(json?: any) {
		if (json) {
			this.id = json.id;
			this.label = json.label;
			this.title = json.title;
			this.metadesc = json.metadesc;
			this.url = json.url;
			this.order = json.order;
			this.locale = json.locale;
			this.code = json.code;
		}
	}

	public toJSON(): any {
		return {
			id: this.id,
			label: this.label,
			title: this.title,
			metadesc: this.metadesc,
			url: this.url,
			order: this.order,
			locale: this.locale,
			code: this.code,
		};
	}
}
