import React, { useContext, useEffect, useState } from 'react';
import {
  ControlledTreeEnvironment,
  Tree,
  TreeItem,
  TreeItemIndex,
} from 'react-complex-tree';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';

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

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);

  const services = ctx.state.activeWorkspace?.spec.services;
  const selectedItem = ctx.state.activeMethod?.fullName;

  const serviceNames = services?.map((it) => it.fullName) || [];
  const [expandedItems, setExpandedItems] =
    useState<TreeItemIndex[]>(serviceNames);

  useEffect(() => {
    setExpandedItems(serviceNames);
  }, [services]);

  const itemsData: Record<TreeItemIndex, TreeItem<Data>> = {
    root: {
      index: 'root',
      isFolder: true,
      data: { display: 'Root item', isMethod: false },
      children: serviceNames,
    },
  };
  services?.forEach((it) => {
    itemsData[it.fullName] = {
      index: it.fullName,
      data: { display: it.displayName, isMethod: false },
      isFolder: true,
      children: it.methods.map((met) => met.fullName),
    };
    it.methods.forEach((met) => {
      itemsData[met.fullName] = {
        index: met.fullName,
        data: { display: met.name, isMethod: true, parent: it.fullName },
      };
    });
  });

  return (
    <div>
      <ControlledTreeEnvironment
        items={itemsData}
        getItemTitle={(item) => item.data.display}
        viewState={{
          '1': {
            expandedItems: expandedItems,
            selectedItems: selectedItem ? [selectedItem] : undefined,
          },
        }}
        onExpandItem={(item) => {
          if (item.data.isMethod) {
            const method = findMethod(
              services,
              item.data.parent || '',
              item.index as string
            );

            ctx.dispatch({ type: 'activeMethod', activeMethod: method! });
          }
          setExpandedItems([...expandedItems, item.index]);
        }}
        onCollapseItem={(item) => {
          setExpandedItems(
            expandedItems.filter(
              (expandedItemIndex) => expandedItemIndex !== item.index
            )
          );
        }}
        onFocusItem={(item) => {
          if (item.data.isMethod) {
            const method = findMethod(
              services,
              item.data.parent || '',
              item.index as string
            );
            ctx.dispatch({ type: 'activeMethod', activeMethod: method! });
          }
        }}
      >
        <Tree treeId="1" rootItem="root" treeLabel="methods collection" />
      </ControlledTreeEnvironment>
    </div>
  );
};
