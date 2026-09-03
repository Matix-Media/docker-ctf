package main

// indexTemplate ist die komplette Weboberflaeche des CTF.
// Sie zeigt immer: Fortschritt, bereits verdiente Flags und genau ein aktuelles Level.
const indexTemplate = `
<!DOCTYPE html>
<html lang="{{.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Brand}}</title>
<style>
  :root {
    --bg: #f4f5f7;
    --card: #ffffff;
    --ink: #1d2330;
    --muted: #5b6472;
    --line: #dfe3e9;
    --blue: #005a9e;
    --green: #1f8a4c;
    --green-bg: #e9f7ef;
    --red: #c62828;
    --red-bg: #fdecec;
    --amber: #b26a00;
    --amber-bg: #fff6e5;
  }
  * { box-sizing: border-box; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background: var(--bg); color: var(--ink);
    max-width: 820px; margin: 0 auto; padding: 24px 16px 64px; line-height: 1.55;
  }
  header { display: flex; align-items: center; justify-content: space-between; gap: 16px; flex-wrap: wrap; }
  header h1 { color: var(--blue); font-size: 1.5rem; margin: 0; }
  .langswitch {
    font-size: .85rem; text-decoration: none; color: var(--blue);
    border: 1px solid var(--line); background: var(--card);
    padding: 6px 12px; border-radius: 999px; white-space: nowrap;
  }
  .langswitch:hover { border-color: var(--blue); }

  section { background: var(--card); border: 1px solid var(--line); border-radius: 10px;
            padding: 18px 20px; margin-top: 18px; }
  h2 { margin: 0 0 4px; font-size: 1.25rem; }
  h3 { margin: 18px 0 4px; font-size: .8rem; text-transform: uppercase;
       letter-spacing: .06em; color: var(--muted); }
  p { margin: 6px 0; }
  code { background: #eceff3; padding: 2px 6px; border-radius: 4px;
         font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: .9em; }
  pre { background: #1d2330; color: #f2f4f7; padding: 14px 16px; border-radius: 8px;
        overflow-x: auto; margin: 10px 0; }
  pre code { background: none; color: inherit; padding: 0; font-size: .85rem; line-height: 1.5; }

  /* Fortschritt */
  .progress { display: flex; align-items: center; gap: 14px; flex-wrap: wrap; }
  .progress .label { font-size: .8rem; text-transform: uppercase; letter-spacing: .06em; color: var(--muted); }
  .dots { display: flex; gap: 8px; list-style: none; margin: 0; padding: 0; }
  .dots li { width: 34px; height: 34px; border-radius: 50%; display: grid; place-items: center;
             font-size: .85rem; font-weight: 600; border: 2px solid var(--line); color: var(--muted); }
  .dots li.done { background: var(--green); border-color: var(--green); color: #fff; }
  .dots li.current { border-color: var(--blue); color: var(--blue); background: #eaf3fa; }

  /* Flags */
  .flag { border-left: 4px solid var(--green); background: var(--green-bg);
          padding: 10px 14px; border-radius: 0 6px 6px 0; margin-bottom: 10px; }
  .flag.fresh { box-shadow: 0 0 0 3px rgba(31,138,76,.25); }
  .flag-head { font-weight: 700; color: var(--green); font-size: .9rem; }
  .flagvalue { display: inline-block; margin-top: 6px; background: #fff;
               border: 1px dashed var(--green); color: var(--green);
               font-weight: 700; letter-spacing: .02em; }
  .note { font-size: .85rem; color: var(--muted); margin-top: 10px; }

  /* Status */
  .status.ok { border-left: 4px solid var(--green); background: var(--green-bg); }
  .status.bad { border-left: 4px solid var(--red); background: var(--red-bg); }
  .status strong { display: block; margin-bottom: 4px; }

  .concept { background: #eaf3fa; border-left: 4px solid var(--blue);
             padding: 10px 14px; border-radius: 0 6px 6px 0; margin-top: 12px; }
  .concept .cap { font-weight: 700; color: var(--blue); font-size: .8rem;
                  text-transform: uppercase; letter-spacing: .06em; }

  .err { border-left: 4px solid var(--red); background: var(--red-bg);
         padding: 10px 14px; border-radius: 0 6px 6px 0; margin-top: 12px; }
  .mount { margin-top: 12px; font-size: .95rem; }
  .check { margin-top: 16px; font-size: .9rem; color: var(--muted);
           border-top: 1px solid var(--line); padding-top: 12px; }
  .check .cap { font-weight: 700; color: var(--ink); display: block; }
  .bonus { margin-top: 12px; background: var(--amber-bg); border-left: 4px solid var(--amber);
           padding: 10px 14px; border-radius: 0 6px 6px 0; font-size: .92rem; }

  details { border: 1px solid var(--line); border-radius: 8px; margin-top: 8px; background: #fafbfc; }
  details[open] { background: #fff; }
  summary { cursor: pointer; padding: 10px 14px; font-weight: 600; font-size: .93rem; }
  summary::marker { color: var(--muted); }
  details .body { padding: 0 14px 14px; }

  form { margin-top: 14px; }
  label { font-weight: 600; font-size: .93rem; display: block; margin-bottom: 6px; }
  input[type=text] { width: 100%; padding: 11px 12px; border: 1px solid var(--line);
                     border-radius: 6px; font-size: 1rem; }
  input[type=text]:focus { outline: 2px solid var(--blue); outline-offset: 1px; border-color: var(--blue); }
  button { margin-top: 10px; width: 100%; padding: 11px; border: 0; border-radius: 6px;
           background: var(--blue); color: #fff; font-size: 1rem; font-weight: 600; cursor: pointer; }
  button:hover { background: #00477d; }

  .cleanup { background: var(--amber-bg); border-color: #f0d9ac; }
  .cleanup .cap { font-weight: 700; color: var(--amber); }
  footer { margin-top: 22px; font-size: .85rem; color: var(--muted); text-align: center; }
  footer p { margin: 4px 0; }
</style>
</head>
<body>

<header>
  <h1>{{.Brand}}</h1>
  <a class="langswitch" href="/?lang={{.OtherLang}}">{{.LangSwitch}}</a>
</header>

<section class="progress">
  <span class="label">{{.ProgressWord}}</span>
  <ol class="dots">
    {{range .Dots}}<li class="{{.State}}">{{.Number}}</li>{{end}}
  </ol>
</section>

{{if .Flags}}
<section>
  {{range .Flags}}
  <div class="flag{{if .Fresh}} fresh{{end}}">
    <div class="flag-head">Level {{.Level}} — {{.Title}}</div>
    <p>{{.Text}}</p>
    <code class="flagvalue">{{.Value}}</code>
  </div>
  {{end}}
  <p class="note">{{.FlagNote}}</p>
</section>
{{end}}

<section class="status {{if .Connected}}ok{{else}}bad{{end}}">
  {{if .Connected}}
    <strong>✅ data-provider-svc</strong>
    <div>{{.ConnOK}}</div>
  {{else}}
    <strong>❌ data-provider-svc</strong>
    <div>{{.ConnBad}}</div>
    <div>{{.Diagnosis}}</div>
  {{end}}
</section>

<section>
  <h2>{{.LevelWord}} {{.Level}} — {{.LevelTitle}}</h2>

  <h3>{{.GoalLabel}}</h3>
  <p>{{.Goal}}</p>

  <div class="concept">
    <span class="cap">{{.ConceptLabel}}</span>
    <p>{{.Concept}}</p>
  </div>

  {{if .Error}}<div class="err">{{.Error}}</div>{{end}}
  {{if .ShowForm4}}<div class="mount">{{.MountNote}}</div>{{end}}

  {{if .ShowForm2}}
  <form method="POST" action="/">
    <input type="hidden" name="step" value="2">
    <label for="answer">{{.Question}}</label>
    <input type="text" id="answer" name="answer" placeholder="{{.Placeholder2}}" autocomplete="off" autofocus>
    <button type="submit">{{.Submit2}}</button>
  </form>
  {{end}}

  {{if .ShowForm4}}
  <form method="POST" action="/">
    <input type="hidden" name="step" value="4">
    <label for="password">{{.PwLabel}}</label>
    <input type="text" id="password" name="password" placeholder="{{.Placeholder4}}" autocomplete="off">
    <button type="submit">{{.Submit4}}</button>
  </form>
  {{end}}

  <h3>{{.HintsLabel}}</h3>
  <details><summary>{{.Hint1Label}}</summary><div class="body">{{.Hint1}}</div></details>
  <details><summary>{{.Hint2Label}}</summary><div class="body">{{.Hint2}}</div></details>
  <details><summary>{{.Hint3Label}}</summary><div class="body">{{.Hint3}}</div></details>

  <div class="check">
    <span class="cap">{{.CheckLabel}}</span>
    {{.Check}}
  </div>

  {{if .Bonus}}<div class="bonus">{{.Bonus}}</div>{{end}}
</section>

<section class="cleanup">
  <span class="cap">{{.CleanupTitle}}</span>
  <p>{{.CleanupBody}}</p>
  <pre><code>docker stop ctf-main &amp;&amp; docker rm ctf-main</code></pre>
</section>

<footer>
  <p>{{.WorkbookNote}}</p>
  <p>{{.StuckNote}}</p>
</footer>

</body>
</html>
`
