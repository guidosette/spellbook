import {Spellbook} from './spellbook';

export class ResponseError {
	Error: string;
	Field: string;

	constructor(error: string, field: string) {
		this.Error = error;
		this.Field = field;
	}
}
