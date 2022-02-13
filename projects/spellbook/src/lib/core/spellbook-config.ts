import {Definition} from './typedef.service';

export class SpellbookConfig {
	apiUrl: string;
	superUserRedirectUrl: string;
	googleMapKey: string;
	typeDefinitions: () => Array<Definition<any>>;
}
