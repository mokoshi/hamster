import { Box, Heading } from '@chakra-ui/react';
import React from 'react';

interface Props {
  title?: string;
}

const Card: React.FC<Props> = (props) => {
  const { title, children } = props;
  return (
    <Box p={2} borderWidth='1px' borderRadius='lg'>
      {title && (
        <Heading m={2} mb={4} size='sm'>
          {title}
        </Heading>
      )}
      {children}
    </Box>
  );
};
export default Card;
