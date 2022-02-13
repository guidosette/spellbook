export class User {

	static readonly PASSWORD_MIN_LEN: number = 8;

	static readonly PERMISSION_ENABLED: string = 'PERMISSION_ENABLED';
	static readonly PERMISSION_EDIT_PERMISSIONS: string = 'PERMISSION_EDIT_PERMISSIONS';
	static readonly PERMISSION_READ_USER: string = 'PERMISSION_READ_USER';
	static readonly PERMISSION_WRITE_USER: string = 'PERMISSION_WRITE_USER';
	static readonly PERMISSION_READ_CONTENT: string = 'PERMISSION_READ_CONTENT';
	static readonly PERMISSION_WRITE_CONTENT: string = 'PERMISSION_WRITE_CONTENT';
	static readonly PERMISSION_READ_MAILMESSAGE: string = 'PERMISSION_READ_MAILMESSAGE';
	static readonly PERMISSION_WRITE_MAILMESSAGE: string = 'PERMISSION_WRITE_MAILMESSAGE';
	static readonly PERMISSION_READ_PLACE: string = 'PERMISSION_READ_PLACE';
	static readonly PERMISSION_WRITE_PLACE: string = 'PERMISSION_WRITE_PLACE';
	static readonly PERMISSION_READ_MEDIA: string = 'PERMISSION_READ_MEDIA';
	static readonly PERMISSION_WRITE_MEDIA: string = 'PERMISSION_WRITE_MEDIA';
	static readonly PERMISSION_READ_PAGE: string = 'PERMISSION_READ_PAGE';
	static readonly PERMISSION_WRITE_PAGE: string = 'PERMISSION_WRITE_PAGE';
	static readonly PERMISSION_READ_SUBSCRIPTION: string = 'PERMISSION_READ_SUBSCRIPTION';
	static readonly PERMISSION_WRITE_SUBSCRIPTION: string = 'PERMISSION_WRITE_SUBSCRIPTION';
	static readonly PERMISSION_READ_ACTION: string = 'PERMISSION_READ_ACTION';
	static readonly PERMISSION_WRITE_ACTION: string = 'PERMISSION_WRITE_ACTION';

	username: string;

	password: string;

	name: string;

	surname: string;

	email: string;

	permissions: Array<string>;

	constructor(json?: any) {
		this.permissions = new Array<string>();
		if (json) {
			this.username = json.username;
			this.name = json.name;
			this.surname = json.surname;
			this.email = json.email;
			if (json.permissions) {
				for (const p of json.permissions) {
					this.permissions.push(p);
				}
			}
		}
	}

	public hasPermission(permission: string): boolean {
		return this.permissions.indexOf(permission) > -1;
	}
}
