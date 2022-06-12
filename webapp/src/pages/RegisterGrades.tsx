import { Button, IconButton, Typography } from '@mui/material';
import { DefaultLayout } from '../components/layout';
import './register-grades.css';
import React, { useState } from 'react';
import PrivacyPolicy from '../components/forms/PrivacyPolicy';
import ContentPasteIcon from '@mui/icons-material/ContentPaste';
import RegisterGrades41 from '../assets/register-grades-4-1.png';
import RegisterGrades42 from '../assets/register-grades-4-2.png';

function RegisterGrades(): JSX.Element {
  const [agree, setAgree] = useState(false);
  const [registerToken, setRegisterToken] = useState<string | null>(null);
  if (!agree) {
    return <PrivacyPolicy onAgree={() => setAgree(true)} />;
  }
  return (
    <DefaultLayout>
      <Typography variant='h4'>成績情報登録</Typography>
      <div className='section'>
        <Typography variant='h6'>
          1. TamperMonkey のインストールをする
        </Typography>
        <Typography variant='body2'>
          ご利用のブラウザをクリックしてインストール作業を行なって下さい。
          <br />
          <a
            target='_blank'
            href='https://chrome.google.com/webstore/detail/tampermonkey/dhdgffkkebhmkfjojejmpbldmpobfkfo?hl=ja'
            rel='noreferrer'
          >
            Chrome
          </a>
          ,{' '}
          <a
            target='_blank'
            href='https://addons.mozilla.org/ja/firefox/addon/tampermonkey/'
            rel='noreferrer'
          >
            Firefox
          </a>
          ,{' '}
          <a
            target='_blank'
            href='https://microsoftedge.microsoft.com/addons/detail/tampermonkey/iikmkjmpaadaobahmlepeloendndfphd?hl=ja-JP'
            rel='noreferrer'
          >
            Microsoft Edge
          </a>
        </Typography>
      </div>

      <div className='section'>
        <Typography variant='h6'>2. UserScript のインストールをする</Typography>
        <Typography variant='body2'>
          <a
            target='_blank'
            href='https://github.com/earlgray283/lab-assignment-system/raw/main/userscript/grades-sender.user.js'
            rel='noreferrer'
          >
            UserScript のリンク
          </a>
        </Typography>
      </div>

      <div className='section'>
        <Typography variant='h6'>3. 成績登録トークンの取得をする</Typography>
        <Button
          variant='contained'
          sx={{ fontSize: '10px' }}
          onClick={() => setRegisterToken('aaaaaaaaaaaa')}
        >
          トークンの取得
        </Button>
        {registerToken && (
          <div
            style={{
              margin: '5px 0px 5px 0px',
              padding: '5px 10px 5px 10px',
              borderRadius: '10px',
              backgroundColor: '#DDDDDD',
              width: 'auto',
            }}
          >
            {registerToken}
            <IconButton size='small' sx={{ alignSelf: 'right' }}>
              <ContentPasteIcon />
            </IconButton>
          </div>
        )}
      </div>

      <div className='section'>
        <Typography variant='h6'>4. 成績情報の送信をする</Typography>
        <Typography variant='body2'>
          「成績情報を送信する」ボタンをクリック
          <br />
          <img src={RegisterGrades41} width='200px' />
          <br />
          先ほど取得したトークンを入力(コピー&ペーストをおすすめします)
          <br />
          <img src={RegisterGrades42} width='200px' />
        </Typography>
      </div>
    </DefaultLayout>
  );
}

export default RegisterGrades;
