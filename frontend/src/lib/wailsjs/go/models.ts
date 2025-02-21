export namespace fingerprint {
	
	export class Config {
	    timezone: string;
	    "webGl:renderer": string;
	    "webGl:vendor": string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timezone = source["timezone"];
	        this["webGl:renderer"] = source["webGl:renderer"];
	        this["webGl:vendor"] = source["webGl:vendor"];
	    }
	}
	export class Fingerprint {
	    version: number;
	    ipv4: string;
	    config: Config;
	    ipv6: string;
	    os: string;
	    key: number;
	    id: string;
	    created_at: string;
	    accessed_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Fingerprint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.ipv4 = source["ipv4"];
	        this.config = this.convertValues(source["config"], Config);
	        this.ipv6 = source["ipv6"];
	        this.os = source["os"];
	        this.key = source["key"];
	        this.id = source["id"];
	        this.created_at = source["created_at"];
	        this.accessed_at = source["accessed_at"];
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
	export class Record {
	    fingerprint: Fingerprint;
	    country: string;
	    proxy_exit_host: string;
	
	    static createFrom(source: any = {}) {
	        return new Record(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fingerprint = this.convertValues(source["fingerprint"], Fingerprint);
	        this.country = source["country"];
	        this.proxy_exit_host = source["proxy_exit_host"];
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
	
	export class Config {
	    auth_code: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auth_code = source["auth_code"];
	    }
	}

}

export namespace profile {
	
	export class Request {
	    name: string;
	    access_code: number;
	    proxy: proxy.Config;
	    imported: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.access_code = source["access_code"];
	        this.proxy = this.convertValues(source["proxy"], proxy.Config);
	        this.imported = source["imported"];
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
	export class Config {
	    id: number;
	    fingerprint: fingerprint.Record;
	    request: Request;
	    hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.fingerprint = this.convertValues(source["fingerprint"], fingerprint.Record);
	        this.request = this.convertValues(source["request"], Request);
	        this.hidden = source["hidden"];
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

export namespace proxy {
	
	export class Config {
	    Host: string;
	    User: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Host = source["Host"];
	        this.User = source["User"];
	        this.Password = source["Password"];
	    }
	}

}

