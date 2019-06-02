package main

/*func StartTasks(c *client.Client, wc *base.Client) {
	fmt.Println(fmt.Sprintf("%+v", wc))
	scheduler.ScheduleTask(scheduler.Task{
		Name:            "Flush sites names as meta",
		FirstCallOnInit: true,
		Period:          5 * time.Minute,
		Fn: func() {
			sites, err := c.Sites()
			if err != nil {
				log.WithError(err).Error("cannot flush metas on sites")
				//return
				sites = []client.Site{
					client.Site{
						ID:          "785627895678903",
						Name:        "default",
						Description: "test",
					},
				}
			}
			sitesB, _ := json.Marshal(sites)

			b := &bytes.Buffer{}
			tpl := template.Must(template.New("").Parse(`
				"{{ .sites }}" JSON-> 'sites' STORE

				[ {{ .rToken }} '~.*' { 'siteId' '~.*' } ] FIND
				<%
					DUP LABELS 'siteId' GET 'site' STORE
					<% $sites $site CONTAINSKEY %>
						<% { $sites $site GET 'siteName' } SETATTRIBUTES  %>
					IFT
				%> FOREACH

			`))

			err = tpl.Execute(b, map[string]interface{}{
				"sites":  string(sitesB),
				"rToken": wc.ReadToken,
				"wToken": wc.WriteToken,
			})
			if err != nil {
				log.WithError(err).Error("cannot prepare sites metas")
			}
			fmt.Println(b.String())

			_res, err := wc.Exec(b.String())
			if err != nil {
				log.WithError(err).Error("cannot flush site metas")
			}
			fmt.Println(string(_res))
		},
	})

}*/
