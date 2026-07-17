export namespace advancement {
	
	export class ItemWithCriteria {
	    ID: string;
	    Title: string;
	    Description: string;
	    Icon: string;
	    Difficulty: string;
	    Branch: string;
	    Done: boolean;
	    IsBig: boolean;
	    Criteria: Record<string, boolean>;
	
	    static createFrom(source: any = {}) {
	        return new ItemWithCriteria(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.Description = source["Description"];
	        this.Icon = source["Icon"];
	        this.Difficulty = source["Difficulty"];
	        this.Branch = source["Branch"];
	        this.Done = source["Done"];
	        this.IsBig = source["IsBig"];
	        this.Criteria = source["Criteria"];
	    }
	}
	export class BranchSnapshot {
	    ID: string;
	    Title: string;
	    Items: ItemWithCriteria[];
	    DoneCount: number;
	    TotalCount: number;
	
	    static createFrom(source: any = {}) {
	        return new BranchSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.Items = this.convertValues(source["Items"], ItemWithCriteria);
	        this.DoneCount = source["DoneCount"];
	        this.TotalCount = source["TotalCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class FullWorldState {
	    WorldName: string;
	    Branches: BranchSnapshot[];
	
	    static createFrom(source: any = {}) {
	        return new FullWorldState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WorldName = source["WorldName"];
	        this.Branches = this.convertValues(source["Branches"], BranchSnapshot);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

export namespace main {
	
	export class AppConfig {
	    minecraftPath: string;
	    pollInterval: number;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minecraftPath = source["minecraftPath"];
	        this.pollInterval = source["pollInterval"];
	    }
	}

}

