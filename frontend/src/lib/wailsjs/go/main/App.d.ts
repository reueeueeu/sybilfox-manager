// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {profile} from '../models';
import {proxy} from '../models';

export function Auth(arg1:number):Promise<void>;

export function GetConfig():Promise<main.Config>;

export function Greet(arg1:string):Promise<string>;

export function HideConfig(arg1:number):Promise<void>;

export function InstallBrowser():Promise<void>;

export function IsAuthed():Promise<boolean>;

export function IsBrowserInstalled():Promise<boolean>;

export function ListConfigs():Promise<Array<profile.Config>>;

export function Logout():Promise<void>;

export function NewConfig(arg1:string,arg2:proxy.Config):Promise<profile.Config>;

export function RunConfig(arg1:number):Promise<void>;

export function SaveConfig(arg1:main.Config):Promise<void>;

export function UpdateConfig(arg1:number,arg2:profile.Request):Promise<profile.Config>;
