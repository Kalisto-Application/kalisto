import React, { useContext } from 'react';
import { Context } from '../state';

import Text from '../ui/Text';

export const VarsError: React.FC = () => {
  const ctx = useContext(Context);

  return (
    <Text text={ctx.state.varsError} props={{ className: 'text-red' }}></Text>
  );
};

export default VarsError;
