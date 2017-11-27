FROM iron/base
EXPOSE 6768

ADD tweetservice-linux-amd64 /
ADD healthchecker-linux-amd64 /
	
HEALTHCHECK --interval=5s --timeout=3s CMD ["./healthchecker-linux-amd64", "-port=6768"] || exit 1
ENTRYPOINT ["./tweetservice-linux-amd64", "-configServerUrl=http://configserver:8888", "-profile=test", "-configBranch=P8"]
