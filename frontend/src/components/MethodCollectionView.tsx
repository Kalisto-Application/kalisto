import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import expandIcon from '../../assets/icons/expand.svg';
import addIcon from '../../assets/icons/plus.svg';
import { CreateRequestFile } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import InputAddItem from '../ui/InputAddItem';
import RequestList from '../ui/RequestList';
const findMethod = (
  s: models.Service[] = [],
  serviceName: string,
  name: string
): models.Method | undefined => {
  for (const service of s) {
    if (service.fullName == serviceName) {
      for (const method of service.methods) {
        if (method.fullName == name) {
          return method;
        }
      }
    }
  }
};

type Data = {
  display: string;
  isMethod: boolean;
  parent?: string;
};

export type dataRequestList = {
  name: string;
  id: string;
};

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);
  const [requestId, setRequestId] = useState('');

  const workspaceID = ctx.state.activeWorkspace?.id;
  const services = ctx.state.activeWorkspace?.spec.services;
  const selectedItem = ctx.state.activeMethod?.fullName;

  const serviceNames = services?.map((it) => it.fullName) || [];
  // debugger
  console.log(services);

  // create tree
  const newServicesName = {
    name: '',
    children:
      services?.map((it, indx) => {
        return {
          name: it.displayName,
          idS: it.fullName,
          children:
            it.methods.map((met) => {
              return {
                name: met.name,
                children: [{ name: '', idM: met.fullName }],
              };
            }) || [],
        };
      }) || [],
  };
  // list request
  const files: dataRequestList[] = [
    { name: 'one', id: '1' },
    { name: 'one', id: '2' },
  ];

  const data = flattenTree(newServicesName);

  const setActiveRequest = (id: string) => {
    setRequestId(id);
  };

  // add Item
  const addItem = (value: string) => {
    // CreateRequestFile(workspaceID, method.FullName, value, content, headers);
  };

  const items = ctx.state.activeWorkspace?.scriptFiles.map((it) => {
    return {
      file: it,
      // inEdit: it.id === isOpenEditInput,
      // isActive: it.id === activeScript?.id,
      onClick: () => setActiveRequest(it.id),
      menu: [
        {
          icon: editIcon,
          text: 'Edit',
          onClick: () => {
            // setIsOpenEditInput(it.id);
          },
        },

        {
          icon: deleteIcon,
          text: 'Delete',
          onClick: () => {
            // setIsOpenDeletePopup(it.id);
          },
        },
      ],
    };
  });

  const [servesicID, setServesicID] = useState('');
  const [NameMethod, setNameMethod] = useState('');

  console.log('servesicID', servesicID);

  return (
    <TreeView
      data={data}
      className="pl-4"
      onExpand={() => {}}
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isExpanded,
        handleExpand,
      }: INodeRendererProps) => {
        {
          return (
            <div
              {...getNodeProps({})}
              style={{
                paddingLeft: 30 * (level <= 2 ? level - 1 : 1.3),
                paddingRight: '30px',
                marginBottom: '10px',
                cursor: 'pointer',
              }}
            >
              <div className="flex ">
                {level <= 2 && <ArrowIcon isOpen={isExpanded} />}
                {element.name}
              </div>
              {level === 3 && (
                <div>
                  <InputAddItem
                    addItem={(value) => addItem(value)}
                    text="New Request"
                    placeholder="Name request"
                  />
                  <div
                    style={{
                      paddingLeft: 15 * (level - 1),
                    }}
                  >
                    <RequestList
                      setActiveRequest={setActiveRequest}
                      files={files}
                      requestId={requestId}
                    />
                    {/* <FileList /> */}
                  </div>
                </div>
              )}
            </div>
          );
        }
      }}
    />
  );
};

/*
Создать файл
CreateRequestFile(workspaceID, method.FullName, name, content, headers)
Обновить созданный файл
UpdateRequestFile(workspaceID, method.FullName, models.File)
Удалить файл по ID
RemoveRequestFile(workspaceID, method.FullName, file.ID)
1) заменить либу и понять что работает ✅
2) научиться отображать фиксированный список в FileList ✅
3) науиться создавать файлы ✅
4) отоброжать активный файл на клик
5) подключить меню
*/
type propsArrowIcon = {
  isOpen: boolean;
};
const ArrowIcon: React.FC<propsArrowIcon> = ({ isOpen }) =>
  isOpen ? (
    <img
      src={expandIcon}
      className="rotate-270 mr-2 transition duration-300 ease-in-out"
    />
  ) : (
    <img
      src={expandIcon}
      className="mr-2 -rotate-90 transition duration-300 ease-in-out"
    />
  );
