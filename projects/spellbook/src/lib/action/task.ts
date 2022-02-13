export class Task {

	name: string;
	url: string;
	scheduleTime: string;
	responseTime: string;
	message: string;
	dispatchCount: number;
	responseCount: number;

	method: string;

	constructor(json?: any) {

		if (json) {
			this.name = json.name;
			this.url = json.url;
			this.scheduleTime = json.scheduleTime;
			this.responseTime = json.responseTime;
			this.message = json.message;
			this.dispatchCount = json.dispatchCount;
			this.responseCount = json.responseCount;

			this.method = json.method;
		}
	}

	public toJSON(): any {
		return {
			name: this.name,
			url: this.url,
			scheduleTime: this.scheduleTime,
			responseTime: this.responseTime,
			message: this.message,
			dispatchCount: this.dispatchCount,
			responseCount: this.responseCount,
			method: this.method,
		};
	}

	public getSimpleName(): string {
		const urls = this.name.split('/');
		if (urls.length > 0) {
			return urls[urls.length - 1];
		}
		return this.url;
	}
}

