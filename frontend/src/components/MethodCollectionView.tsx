import React, { useState, useEffect } from "react";
import { ControlledTreeEnvironment, Tree, TreeItemIndex, TreeItem } from 'react-complex-tree';
import "react-complex-tree/lib/style-modern.css";

interface MethodCollectionProps {
  setActiveMethod: (fullName: MethodItem) => void;
  services: ServiceItem[];
  selectedItem?: string;
};

export type ServiceItem = {
  name: string;
  fullName: string;
  methods: MethodItem[];
}

const findMethod = (s: ServiceItem[], name: string): MethodItem | undefined => {
  for (const service of s) {
    for (const method of service.methods) {
      if (method.fullName == name) {
        return method
      }
    }
  }
}

export type MethodItem = {
  name: string;
  fullName: string;
  requestExample: string
} 

type Data = {
  display: string;
  isMethod: boolean;
}

export const MethodCollection: React.FC<MethodCollectionProps> = ({ setActiveMethod, services, selectedItem }) => {
  const serviceNames = services.map(it => it.fullName)
  const [expandedItems, setExpandedItems] = useState<TreeItemIndex[]>(serviceNames);

  useEffect(() => {
    setExpandedItems(serviceNames)
  }, [services])

  const itemsData: Record<TreeItemIndex, TreeItem<Data>> = {root: {
    index: 'root',
    isFolder: true,
    data: {display: 'Root item', isMethod: false},
    children: serviceNames,
  }}
  services.forEach(it => {
    itemsData[it.fullName] = {
      index: it.fullName,
      data: {display: it.name, isMethod: false},
      isFolder: true,
      children: it.methods.map(met => met.fullName),
    }
    it.methods.forEach(met => {
      itemsData[met.fullName] = {
        index: met.fullName,
        data: {display: met.name, isMethod: true},
      }
    })
  })

  return(
    <div>
    <ControlledTreeEnvironment
      items={itemsData}
      getItemTitle={item => item.data.display}
      viewState={{
        "1": {
          expandedItems: expandedItems,
          focusedItem: selectedItem,
        }
      }}
      onExpandItem={item => {
        if (item.data.isMethod) {
          const method = findMethod(services, item.index as string)
          setActiveMethod(method!);
        }
        setExpandedItems([...expandedItems, item.index])}
      }
      onCollapseItem={item => {
        setExpandedItems(expandedItems.filter(expandedItemIndex => expandedItemIndex !== item.index))
      }}
      onFocusItem={item => {
        if (item.data.isMethod) {
          const method = findMethod(services, item.index as string)
          setActiveMethod(method!);
        }
      }}
    >
        <Tree treeId="1" rootItem="root" treeLabel="methods collection" />
    </ControlledTreeEnvironment>
    </div>
    );
};
