import { IconButton, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import './copy-box.css';
import ContentPasteIcon from '@mui/icons-material/ContentPaste';

function CopyBox(props: {
  children: React.ReactNode;
  copyText: string;
}): JSX.Element {
  const [openTooltip, setOpenTooltip] = useState(false);
  return (
    <div className='copy-box'>
      {props.children}
      <Tooltip
        arrow
        open={openTooltip}
        onClose={() => setOpenTooltip(false)}
        placement='top'
        title='Copied'
      >
        <IconButton size='small' sx={{ alignSelf: 'right' }}>
          <ContentPasteIcon />
        </IconButton>
      </Tooltip>
    </div>
  );
}

export default CopyBox;
