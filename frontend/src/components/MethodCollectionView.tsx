import React, { useContext, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import expandIcon from '../../assets/icons/expand.svg';
import gIcon from '../../assets/icons/g.svg';

import {
  CreateRequestFile,
  RemoveRequestFile,
  UpdateRequestFile,
} from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import CreateItem from '../ui/CreateItem';

import FileList from '../ui/FileList';
import DeletePopup from './DeletePopup';

const findMethod = (
  s: models.Service[] = [],
  serviceName: string,
  metName: string
): models.Method | undefined => {
  for (const service of s) {
    if (service.fullName == serviceName) {
      for (const method of service.methods) {
        if (method.fullName == metName) {
          return method;
        }
      }
    }
  }
};

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState('');

  if (!ctx.state.activeWorkspace) {
    return <></>;
  }

  const requestFiles = ctx.state.activeWorkspace?.requestFiles
    ? ctx.state.activeWorkspace?.requestFiles
    : {};
  const workspaceID = ctx.state.activeWorkspace?.id || '';
  const services = ctx.state.activeWorkspace?.spec.services;
  const activeRequestID = ctx.state.activeRequestFileId;

  // creatign a data tree
  const newServicesName = {
    name: '',
    children:
      services?.map((it) => {
        return {
          name: it.displayName,

          children:
            it.methods.map((met) => {
              return {
                name: met.name,
                children: [{ name: met.fullName }],
              };
            }) || [],
        };
      }) || [],
  };
  const data = flattenTree(newServicesName);

  //  active request and activeMethod
  const setActiveRequestMethod = (id: string, fullNameMet: string) => {
    ctx.dispatch({ type: 'setActiveRequest', id });

    const servesActive = services.find((it) => {
      if (fullNameMet.startsWith(it.fullName)) return it;
    });

    if (!servesActive) return;

    const method = findMethod(services, servesActive.fullName, fullNameMet);

    ctx.dispatch({ type: 'activeMethod', activeMethod: method! });
  };

  // add request
  const addItem = (value: string, fullNameMet: string) => {
    const servesActive = services.find((it) => {
      if (fullNameMet.startsWith(it.fullName)) return it;
    });
    if (!servesActive) return;

    const met = findMethod(services, servesActive.fullName, fullNameMet);
    if (!met) return;

    CreateRequestFile(
      workspaceID,
      fullNameMet,
      value,
      met?.requestExample,
      ''
    ).then((res) => {
      ctx.dispatch({
        type: 'addRequestFile',
        file: {
          ...requestFiles,
          [fullNameMet]: [...(requestFiles[fullNameMet] || []), res],
        },
      });
    });
  };

  // delete request
  const deleteRequest = (metName: string) => {
    if (!isOpenDeletePopup) return;

    RemoveRequestFile(workspaceID, metName, isOpenDeletePopup).then((res) => {
      let ws = new models.Workspace({
        ...ctx.state.activeWorkspace,
        requestFiles: {
          [metName]: [...(res[metName] || [])],
        },
      });

      ctx.dispatch({ type: 'updateWorkspace', workspace: ws });
    });
  };
  // Edit request
  const renameRequest = (name: string, metName: string) => {
    const renamed = new models.File({
      ...ctx.state.activeWorkspace?.requestFiles[metName].find(
        (it) => it.id === isOpenEditInput
      ),
      name: name,
    });

    UpdateRequestFile(workspaceID, metName, renamed).then((res) => {
      ctx.dispatch({
        type: 'updateRequestFile',
        file: renamed,
        metName,
      });
    });
  };

  const paddingSteps = [0, 8, 16];

  return (
    <TreeView
      data={data}
      defaultExpandedIds={data[0].children}
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isExpanded,
        isBranch,
      }: INodeRendererProps) => {
        {
          return (
            <div
              onKeyDown={(e) => e.stopPropagation()}
              {...getNodeProps({})}
              style={{
                cursor: 'pointer',
                width: '100%',
                height: '32px',
                display: 'inline-block',
              }}
            >
              {level <= 2 && (
                <div
                  style={{ paddingLeft: ` ${paddingSteps[level]}px` }}
                  className="flex w-full hover:bg-borderFill"
                >
                  {level <= 2 && <ArrowIcon isOpen={isExpanded} />}
                  {level <= 2 && element.name}
                </div>
              )}
              {level === 3 && (
                <div className="w-full flex-[0_1_100%] flex-col pl-[22px]">
                  <CreateItem
                    fullNameMet={element.name}
                    addItem={(value, fullNameMet) =>
                      addItem(value, fullNameMet || '')
                    }
                    text="New Request"
                    placeholder="Name request"
                  />

                  <FileList
                    items={(requestFiles[element.name] || []).map((it) => {
                      return {
                        file: it,
                        inEdit: it.id === isOpenEditInput,
                        isActive: it.id === activeRequestID,
                        onClick: () =>
                          setActiveRequestMethod(it.id, element.name),
                        menu: [
                          {
                            icon: editIcon,
                            text: 'Edit',
                            onClick: () => {
                              setIsOpenEditInput(it.id);
                            },
                          },
                          {
                            icon: deleteIcon,
                            text: 'Delete',
                            onClick: () => {
                              setIsOpenDeletePopup(it.id);
                            },
                          },
                        ],
                      };
                    })}
                    onCloseInput={() => setIsOpenEditInput('')}
                    editFile={(id: string) => renameRequest(id, element.name)}
                    gIcon={gIcon}
                  />
                  <DeletePopup
                    id={isOpenDeletePopup}
                    isOpen={isOpenDeletePopup !== ''}
                    onClose={() => setIsOpenDeletePopup('')}
                    deleteScript={() => {
                      deleteRequest(element.name);
                    }}
                    title="Delete script?"
                  />
                </div>
              )}
            </div>
          );
        }
      }}
    />
  );
};

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
