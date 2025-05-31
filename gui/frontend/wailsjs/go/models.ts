export namespace main {
	
	export class Change {
	    type: string;
	    target: string;
	    newName: string;
	    selected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Change(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.target = source["target"];
	        this.newName = source["newName"];
	        this.selected = source["selected"];
	    }
	}

}

