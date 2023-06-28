import React, {useState} from "react";
import { UrlInput } from "../components/UrlInput";
import { CodeEditor } from "../components/CodeEditor";
import { MethodCollection, MethodItem } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"

type ContentProps = {
    workspaceId: string;
    methodItems: MethodItem[];
  }

export const MainPageContent: React.FC<ContentProps> = ({workspaceId, methodItems}) => {
    const [inputText, setInputText] = useState<string>('');
    const [url, setUrl] = useState<string>('localhost:9000');
  
    const [outText, setOutText] = useState<string>('');
  
    const sendRequest = (event: React.SyntheticEvent) => {
      const body = `{id: "yoba"}`
      SendGrpc({addr: url, workspaceId: workspaceId, method: methodItems[0].fullName, body: body, meta: ""}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };
  
    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection onClick={()=>{}} items={methodItems} />
        <CodeEditor text={inputText} setText={setInputText}/>
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  