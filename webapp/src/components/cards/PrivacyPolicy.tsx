import { Box, Typography } from '@mui/material';
import React from 'react';

export function PrivacyPolicy(): JSX.Element {
  return (
    <Box>
      <Typography variant='h4'>1. プライバシーポリシー</Typography>
      <Typography variant='h6'>
        1.1. 収集する成績情報の使用用途について
      </Typography>
      <Typography variant='body2'>
        <ul>
          <li>収集するデータは「科目名」「単位数」「点数」「報告日」です。</li>
          <li>
            2022年研究室配属用 GPA を計算することを目的として使用されます。
          </li>
          <li>
            成績情報は一時的にサーバーに送信されますが、GPA
            の計算が完了した時点で破棄されるため、サーバー上に保存されることはありません。
          </li>
        </ul>
      </Typography>
      <Typography variant='h6'>1.2. 2022年研究室配属用 GPA について</Typography>
      <Typography variant='body2'>
        <ul>
          <li>
            2022年研究室配属用 GPA では、<strong>点数が60点以上</strong> かつ{' '}
            <strong>報告日が2022年3月31日</strong>{' '}
            までである科目のみが計算に使用されます。すなわち、成績が不可の科目は、研究室配属のための
            GPA に反映されません。
          </li>
          <li>
            研究室の平均 GPA 情報などを表示するため、GPA
            自体は公開されます。しかし、その GPA
            から個人が特定されることはありません。
          </li>
        </ul>
      </Typography>
    </Box>
  );
}
