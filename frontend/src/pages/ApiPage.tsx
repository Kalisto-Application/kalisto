import React, {useState} from "react";
import { UrlInput } from "../components/UrlInput";
import { MethodCollection } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"
import {models} from "../../wailsjs/go/models"
import { RequestEditor } from "../components/RequestEditor";

type ApiPageProps = {
    workspace: models.Workspace;
    method?: models.Method;
    inputText: string;
}

export const ApiPage: React.FC<ApiPageProps> = ({workspace, method, inputText }) => {
    const [url, setUrl] = useState<string>('localhost:9000');
    const [outText, setOutText] = useState<string>('');
  
    const sendRequest = (event: React.SyntheticEvent) => {
      if (method?.fullName == '') {
        //TODO: disable Send button
        return
      }

      SendGrpc({addr: url, workspaceId: workspace.id, method: method!.fullName, body: inputText, meta: ""}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };

    if (!workspace) {
      return (<div></div>)
    }

    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection services={workspace.spec.services} selectedItem={method?.fullName} />
        <RequestEditor />
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  