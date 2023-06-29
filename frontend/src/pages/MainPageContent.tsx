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
    let firstMethod: string = '';
    if (methodItems.length > 0 && methodItems[0].methods.length > 0) {
      firstMethod = methodItems[0].methods[0].fullName
    }

    const [inputText, setInputText] = useState<string>('');
    const [url, setUrl] = useState<string>('localhost:9000');
    const [method, setMethod] = useState<string>(firstMethod);
  
    const [outText, setOutText] = useState<string>('');
  
    const sendRequest = (event: React.SyntheticEvent) => {
      if (method == '') {
        //TODO: disable Send button
      }

      const body = `{id: "yoba"}`

      SendGrpc({addr: url, workspaceId: workspaceId, method: method, body: body, meta: ""}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };

    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection setActiveMethod={setMethod} items={methodItems} defaultFocused={firstMethod} />
        <CodeEditor text={inputText} setText={setInputText}/>
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  