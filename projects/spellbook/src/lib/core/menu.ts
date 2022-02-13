export class Menu {

	public children: Menu[] = [];

	constructor(public readonly id: string, public readonly url: string, public readonly displayName: string, public readonly icon: string, public readonly queryParams?: any) {
	}

	public addChildren(children: Menu) {
		this.children.push(children);
	}
}
