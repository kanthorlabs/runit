FROM {{ .Version }}

# Set environment variables to prevent Python from writing pyc files and to buffer output
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Create a non-root user
RUN groupadd --gid 1000 kanthorlab \
    && useradd --uid 1000 --gid kanthorlab --create-home kanthorlab

# Set working directory
WORKDIR /app

# Switch to non-root user
USER kanthorlab

# Install Python dependencies
COPY --chown=kanthorlab:kanthorlab requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --chown=kanthorlab:kanthorlab main.py .

{{- range .Ports }}
EXPOSE {{ . }}
{{- end }}

CMD ["python", "main.py"]

