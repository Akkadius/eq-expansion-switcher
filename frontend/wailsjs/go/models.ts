export namespace config {
	
	export class Config {
	    eq_dir: string;
	    current_expansion: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.eq_dir = source["eq_dir"];
	        this.current_expansion = source["current_expansion"];
	    }
	}

}

export namespace eqassets {
	
	export class Expansion {
	    id: number;
	    name: string;
	    mask: number;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new Expansion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.mask = source["mask"];
	        this.icon = source["icon"];
	    }
	}
	export class ExpansionFiles {
	    expansion: Expansion;
	    files: string[];
	
	    static createFrom(source: any = {}) {
	        return new ExpansionFiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.expansion = this.convertValues(source["expansion"], Expansion);
	        this.files = source["files"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

