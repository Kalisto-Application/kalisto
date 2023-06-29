import React, {useState} from "react";
import { UrlInput } from "../components/UrlInput";
import { CodeEditor } from "../components/CodeEditor";
import { MethodCollection, ServiceItem, MethodItem } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"

type ContentProps = {
    workspaceId: string;
    methodItems: ServiceItem[];
  }

export const MainPageContent: React.FC<ContentProps> = ({workspaceId, methodItems}) => {
    var firstMethod: MethodItem = {name: "", fullName: "", requestExample: ""};
    if (methodItems.length > 0 && methodItems[0].methods.length > 0) {
      firstMethod = methodItems[0].methods[0]
    }

    const [inputText, setInputText] = useState<string>(firstMethod.requestExample);
    const [url, setUrl] = useState<string>('localhost:9000');
    const [method, setMethod] = useState<MethodItem>(firstMethod);
  
    const [outText, setOutText] = useState<string>('');
  
    const sendRequest = (event: React.SyntheticEvent) => {
      if (method.fullName == '') {
        //TODO: disable Send button
      }

      SendGrpc({addr: url, workspaceId: workspaceId, method: method.fullName, body: inputText, meta: ""}).then(res => {
        setOutText(res.body)
      }).catch(err => {
        setOutText(err)
      })
    };

    const onSetMethod = (method: MethodItem) => {
      setMethod(method);
      setInputText(method.requestExample)
    }

    return (
      <div className="p-4">
        <UrlInput onClick={sendRequest} value={url} setValue={setUrl} />
        <div className="flex flex-1">
        <MethodCollection setActiveMethod={onSetMethod} items={methodItems} defaultFocused={firstMethod} />
        <CodeEditor text={inputText} setText={setInputText}/>
        <span>{outText}</span>
        </div>
      </div>
    );
  };
  