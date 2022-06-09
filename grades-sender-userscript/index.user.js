// ==UserScript==
// @name         grades-sender
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  成績情報を取得し、lab-assignment-system に送信する。
// @author       earlgray
// @match        https://gakujo.shizuoka.ac.jp/kyoumu/seisekiSearchStudentInit.do*
// @grant        none
// ==/UserScript==

(function () {
  'use strict';

  /**
   *
   * @param {HTMLDivElement} html
   * @returns
   */
  const scrapeStudentInfo = async (html) => {
    const pElem = html.querySelector(
      'div > table:nth-child(10) > tbody > tr:nth-child(2) > td > table > tbody > tr > td > table > tbody > tr > td > p'
    );
    if (!pElem) {
      return;
    }
    const matches = pElem.textContent
      .trim()
      .match(/学籍番号：(\d+)　\n\t\t\t\t  学生氏名：(.+)/i);
    if (matches.length !== 3) {
      alert('matches.length !== 3');
      return;
    }
    const studentNumber = Number(matches[1]);
    const studentName = matches[2];
    return {
      studentNumber: studentNumber,
      studentName: studentName,
    };
  };

  /**
   *
   * @param {HTMLDivElement} html
   * @returns
   */
  const scrapeGradesTable = async (html) => {
    const trList = html.querySelectorAll(
      'div > table:nth-child(12) > tbody > tr > td > table > tbody > tr'
    );
    if (trList === 0) {
      return;
    }

    const grades = new Array(0);
    let unitSum = Number(0);
    let gpSum = Number(0);
    for (const tr of Array.from(trList).slice(1)) {
      const grade = tr.querySelector('td:nth-child(6)').textContent.trim();
      if (grade === '不可' || grade === '合') {
        continue;
      }
      const unitNum = Number(
        tr.querySelector('td:nth-child(5)').textContent.trim()
      );
      const gp = Number(tr.querySelector('td:nth-child(8)').textContent.trim());
      const reportedAt = tr
        .querySelector('td:nth-child(10)')
        .textContent.trim();
      grades.push({ unitNum: unitNum, gp: gp, reportedAt: reportedAt });

      unitSum += unitNum;
      gpSum += gp * unitNum;
    }
    const gpa = gpSum / unitSum;

    return {
      unitSum: unitSum,
      gpSum: gpSum,
      gpa: gpa,
      grades: grades,
    };
  };

  const postJson = async (url, data) => {
    return await fetch(url, {
      method: 'POST',
      mode: 'no-cors',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  };

  (async function () {
    const currentUrl = window.location.href;
    const resp = await fetch(currentUrl);
    if (!resp.ok) {
      console.error(`failed to fetch ${currentUrl}`);
      console.error('Maybe it is a session error');
      return;
    }
    const tmpElem = document.createElement('div');
    tmpElem.innerHTML = await resp.text();
    console.log(tmpElem);
    const gradeObj = await scrapeGradesTable(tmpElem);
    const studentInfo = await scrapeStudentInfo(tmpElem);
    let jsonObj = gradeObj;
    jsonObj.studentName = studentInfo.studentName;
    jsonObj.studentNumber = studentInfo.studentNumber;

    const backendUrl = 'https://hogehoge.com';
    console.log(JSON.stringify(jsonObj));
    //await postJson(backendUrl, jsonObj);
  })().catch((e) => {
    console.error(e);
  });
})();
