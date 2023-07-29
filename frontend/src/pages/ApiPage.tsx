import React, {useContext, useState} from "react";
import { UrlInput } from "../components/UrlInput";
import { MethodCollection } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"
import { RequestEditor } from "../components/RequestEditor";
import { Context } from "../state";

export const ApiPage: React.FC = () => {
    const ctx = useContext(Context);

    const [url, setUrl] = useState<string>('localhost:9000');
    const [outText, setOutText] = useState<string>('');

    const sendRequest = (_: React.SyntheticEvent) => {
      if (!ctx.state.activeMethod || !ctx.state.activeWorkspace) {
        //TODO: disable Send button
        return
      }

      SendGrpc!({addr: url, workspaceId: ctx.state.activeWorkspace.id, method: ctx.state.activeMethod.fullName, body: ctx.state.requestText, meta: ctx.state.requestMetaText}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };

    if (!ctx.state.activeWorkspace) {
      return (<div></div>)
    }

    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection services={ctx.state.activeWorkspace.spec.services} selectedItem={ctx.state.activeMethod?.fullName} />
        <RequestEditor />
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  