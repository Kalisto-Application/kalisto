// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';

export function DeleteEnvivonment(arg1:string,arg2:string):Promise<void>;

export function DeleteWorkspace(arg1:string):Promise<void>;

export function EnvironmentsByWorkspace(arg1:string):Promise<models.Envs>;

export function FindWorkspaces():Promise<Array<models.Workspace>>;

export function GetWorkspace(arg1:string):Promise<models.Workspace>;

export function NewWorkspace():Promise<models.Workspace>;

export function RenameWorkspace(arg1:string,arg2:string):Promise<void>;

export function SaveEnvironment(arg1:models.EnvRaw):Promise<models.Env>;

export function SendGrpc(arg1:models.Request):Promise<models.Response>;
