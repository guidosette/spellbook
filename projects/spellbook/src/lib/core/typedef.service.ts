/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-08-24.
 */
import {Injectable} from '@angular/core';


export interface Definition<T> {
	definitionFor(): string;

	field(id: string): Field;

	menuLabel(): string;

	menuIcon(): string;
}

export class Field {
	type: string;
	public hidden: boolean;

	constructor(public id: string, public label: string) {
		// todo: implement type when dynamic types will be supported
		this.hidden = false;
	}
}

@Injectable()
export class TypedefService {

	private definitions: Map<string, Definition<any>>;

	constructor() {
		this.definitions = new Map<string, any>();
	}

	public addTypeDefinition<T>(definition: Definition<T>) {
		const def = this.definitions.get(definition.definitionFor());
		if (!def) {
			this.definitions.set(definition.definitionFor(), definition);
			return;
		}

		throw new Error(`Definition already given for ${definition.definitionFor()}: ${JSON.stringify(def)}`);
	}

	public getTypeDefinition<T>(type: string): Definition<T> {
		return this.definitions.get(type);
	}
}
